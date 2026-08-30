// Package web is the browser client for the Vactrol card game. It renders an
// interactive two-player hotseat match with the same engine the terminal UI
// uses, compiled to WebAssembly via the go-app framework: the playtester sees
// the whole board and drives both sides.
package web

import (
	"strconv"
	"time"

	"github.com/dmikalova/vactrol/internal/cards"
	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/dmikalova/vactrol/internal/match"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// phase is the interaction state of the client, distinct from the engine's own
// turn structure: it tracks what the UI is currently asking the player to do.
type phase int

const (
	phaseHouse       phase = iota // start of turn: pick the active house
	phaseMain                     // select cards and act
	phaseFlank                    // a creature is selected to play; pick a flank
	phaseFightTarget              // a fight was declared; pick an enemy to hit
	phaseOver                     // the game is finished
)

// selKind labels what the currently selected card is, so the action bar can
// offer the right verbs.
type selKind int

const (
	selNone         selKind = iota
	selHand                 // a card in the active player's hand
	selYourCreature         // a creature the active player controls
	selYourArtifact         // an artifact the active player controls
	selOther                // any other card (read-only: opponent's, an upgrade)
)

// chooseReply carries the player's answer to an engine chooser request: the
// chosen id, or ok=false when the player cancels.
type chooseReply struct {
	id engine.LocalID
	ok bool
}

// webChooser adapts the engine's synchronous Chooser to go-app's single UI
// goroutine. The engine calls ChooseCreature from a background action goroutine
// (see game.runAction); it shows the chooser overlay on the UI goroutine, then
// blocks until a candidate is clicked (or the request is cancelled).
type webChooser struct {
	g           *game
	reply       chan chooseReply
	optionReply chan int
}

// ChooseCreature posts a chooser request to the UI and waits for the player's
// pick. It returns false when there are no candidates or the player cancels.
func (c *webChooser) ChooseCreature(_, prompt string, candidates []engine.LocalID) (engine.LocalID, bool) {
	if len(candidates) == 0 {
		return 0, false
	}
	// Discard any stale reply left in the buffer (e.g. from a double click on the
	// previous prompt) so it cannot silently answer this one.
	select {
	case <-c.reply:
	default:
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosing = true
		c.g.chooserPrompt = prompt
		c.g.chooserCandidates = candidates
	})
	r := <-c.reply
	c.g.dispatch(func(app.Context) {
		c.g.choosing = false
		c.g.chooserPrompt = ""
		c.g.chooserCandidates = nil
	})
	return r.id, r.ok
}

// ChooseOption implements the engine's OptionChooser: it posts a labeled
// multiple-choice prompt (e.g. whether to take archived cards into hand) to the
// UI and blocks until the player clicks one of the option buttons. Without this,
// the engine falls back to the first option — which silently auto-took archives.
func (c *webChooser) ChooseOption(_, prompt string, options []string) int {
	// Drop any stale reply so a leftover click cannot answer this prompt.
	select {
	case <-c.optionReply:
	default:
	}
	c.g.dispatch(func(app.Context) {
		c.g.choosingOption = true
		c.g.optionPrompt = prompt
		c.g.optionLabels = options
	})
	i := <-c.optionReply
	c.g.dispatch(func(app.Context) {
		c.g.choosingOption = false
		c.g.optionPrompt = ""
		c.g.optionLabels = nil
	})
	return i
}

// NewGame returns the root component for a fresh browser client session. The
// match itself is seeded on the client in OnMount.
func NewGame() app.Composer { return &game{selHand: -1, zonesPlayer: -1} }

// game is the root component: it owns the live engine.Game and all UI state.
type game struct {
	app.Compo

	g          *engine.Game
	chooser    *webChooser
	seed       int64             // deal seed; persisted so a hot-reload can rebuild the match
	deckHouses [2][]engine.House // each player's three deck houses (house choices)

	// dispatch schedules a mutation on the UI goroutine (captured from a Context).
	// It lets the background chooser update fields safely.
	dispatch func(func(app.Context))

	phase phase
	busy  bool // an action goroutine is resolving; input is ignored

	// selection
	sel      engine.LocalID
	selKind  selKind
	selHand  int // hand index when selKind == selHand, else -1
	hasSel   bool
	attacker engine.LocalID // creature declared to fight, during phaseFightTarget

	// engine chooser overlay
	choosing          bool
	chooserPrompt     string
	chooserCandidates []engine.LocalID

	// engine option chooser: a labeled multiple choice (e.g. take archives?)
	choosingOption bool
	optionPrompt   string
	optionLabels   []string

	// zonesPlayer, when >= 0, opens the out-of-play zone viewer (discard, archives,
	// and purge piles) for that player. -1 keeps the viewer closed.
	zonesPlayer int

	// Manual (sandbox) mode lives on the engine (g.g.Manual()); these back its UI:
	// the fuzzy card picker's open state and query, and the cached card pool it
	// searches.
	pickerOpen  bool
	pickerQuery string
	allDefs     []engine.CardDefinition

	// dragging is set while a hand card is being dragged, so the board shows as a
	// drop zone. It is cleared when the drag ends or the card is dropped.
	dragging bool

	// detailDef, when non-nil, overrides the detail panel to show a card looked up
	// from the log by name (a printed-text reference, not a live board selection).
	detailDef *engine.CardDefinition
	// defByName maps every card's name to its definition, for log-mention lookups.
	defByName map[string]*engine.CardDefinition

	status string // transient message (usually an action error)
}

// persistKey names the session-storage slot holding the in-progress match, so a
// hot-reload of the wasm can resume it instead of dealing a new game.
const persistKey = "vactrol.match"

// snapshotVersion tags persisted state; bump it when an engine change makes older
// snapshots invalid so a stale one is flushed instead of restored.
const snapshotVersion = 2

// snapshot is the persisted match. The seed deterministically rebuilds the
// catalog and card ids; the flat GameState carries everything mutable. All other
// state (choosers, deck houses, card index, UI phase) is reconstructed from these.
type snapshot struct {
	Version int
	Seed    int64
	State   engine.GameState
}

// OnMount resumes the saved match if there is one, else deals a fresh game. It
// runs on the UI goroutine once the component is inserted into the page.
func (g *game) OnMount(ctx app.Context) {
	g.dispatch = ctx.Dispatch
	if !g.resume(ctx) {
		g.newMatch()
	}
	g.save(ctx)
}

// OnAppUpdate fires when go-app detects a freshly built wasm bundle. It persists
// the current match and reloads the page onto the new build, which OnMount then
// resumes — the hot-reload path that keeps the game state.
func (g *game) OnAppUpdate(ctx app.Context) {
	g.save(ctx)
	ctx.Reload()
}

// save writes the current match to session storage. It runs after every action
// and before a hot-reload, so a reload resumes from the latest state.
func (g *game) save(ctx app.Context) {
	if g.g == nil {
		return
	}
	_ = ctx.SessionStorage().Set(persistKey, snapshot{
		Version: snapshotVersion,
		Seed:    g.seed,
		State:   g.g.State,
	})
}

// resume rebuilds the match from a saved snapshot, reporting whether it restored
// one. A missing, wrong-version, or engine-incompatible snapshot is dropped and
// resume returns false, so the caller deals a fresh game. Reconstruction from the
// seed is deterministic, so the rebuilt catalog and card ids match the saved
// state exactly; any residual mismatch (an older card pool) is caught before use.
func (g *game) resume(ctx app.Context) (ok bool) {
	store := ctx.SessionStorage()
	if !store.Contains(persistKey) {
		return false
	}
	var snap snapshot
	if err := store.Get(persistKey, &snap); err != nil ||
		snap.Version != snapshotVersion || snap.Seed == 0 {
		store.Del(persistKey)
		return false
	}
	// Rebuilding or reading a state from a different engine/card pool can panic on
	// an out-of-range id; recover and fall back to a fresh deal.
	defer func() {
		if recover() != nil {
			store.Del(persistKey)
			ok = false
		}
	}()

	g.seed = snap.Seed
	eg, houses := match.New("Player 1", "Player 2", snap.Seed)
	g.install(eg, houses)
	g.g.State = snap.State
	if !g.stateReadsCleanly() {
		store.Del(persistKey)
		return false
	}
	g.settlePhase()
	g.clearSelection()
	g.zonesPlayer = -1
	g.status = ""
	return true
}

// stateReadsCleanly reports whether every card id in the restored state resolves
// against the rebuilt catalog. A snapshot from a different card pool can hold ids
// out of range, which panics on lookup; recovering turns that into a clean
// start-over signal.
func (g *game) stateReadsCleanly() (ok bool) {
	defer func() {
		if recover() != nil {
			ok = false
		}
	}()
	for p := 0; p < 2; p++ {
		for _, zone := range [][]engine.LocalID{
			g.g.Hand(p), g.g.Deck(p), g.g.Battleline(p), g.g.Artifacts(p),
			g.g.Discard(p), g.g.Archives(p), g.g.Purge(p),
		} {
			for _, id := range zone {
				g.g.Def(id)
			}
		}
	}
	return true
}

// newMatch seeds a new game, wires the shared human chooser to both players, and
// deals random decks. Both sides are driven by the same person (hotseat).
func (g *game) newMatch() {
	g.seed = time.Now().UnixNano()
	eg, houses := match.New("Player 1", "Player 2", g.seed)
	g.install(eg, houses)
	g.g.BeginTurn(0) // no Æmber yet, so the opening forge step is a no-op
	g.phase = phaseHouse
	g.clearSelection()
	g.zonesPlayer = -1
	g.status = ""
}

// install wires a freshly built engine game into the component: it attaches the
// shared human chooser to both players and records the harness and its deck
// houses. The caller sets the starting phase.
func (g *game) install(eg *engine.Game, houses [2][]engine.House) {
	ch := &webChooser{g: g, reply: make(chan chooseReply, 1), optionReply: make(chan int, 1)}
	eg.SetChooser(0, ch)
	eg.SetChooser(1, ch)
	// Declare each player's deck houses so a forced house they lack is ignored
	// (cannot overrides must); the human is prompted for each key's colour.
	eg.SetPlayerHouses(0, houses[0])
	eg.SetPlayerHouses(1, houses[1])
	g.g = eg
	g.chooser = ch
	g.deckHouses = houses
	if g.defByName == nil {
		g.defByName = cardsByName()
	}
	if g.allDefs == nil {
		g.allDefs = cards.All()
	}
}

// cardsByName indexes every registered card's definition by its display name so
// the log can link the cards it mentions. The pointers reference the returned
// slice's backing array, which the map keeps alive.
func cardsByName() map[string]*engine.CardDefinition {
	all := cards.All()
	m := make(map[string]*engine.CardDefinition, len(all))
	for i := range all {
		m[all[i].Name] = &all[i]
	}
	return m
}

// active returns the player whose turn it is; the client always renders from
// this player's perspective (their board and hand at the bottom).
func (g *game) active() int { return g.g.State.ActivePlayer }

func (g *game) clearSelection() {
	g.sel = 0
	g.selKind = selNone
	g.selHand = -1
	g.hasSel = false
	g.attacker = 0
	g.detailDef = nil
}

// ---- action plumbing ----

// runAction resolves an engine mutation on a background goroutine so the UI stays
// responsive while an effect blocks on the chooser. When it finishes it advances
// the phase and clears any transient selection on the UI goroutine.
func (g *game) runAction(ctx app.Context, fn func() error) {
	if g.busy {
		return
	}
	g.busy = true
	g.status = ""
	// Complete via g.dispatch, which is bound in OnMount to the always-mounted
	// root component. Do NOT use ctx.Dispatch here: ctx is tied to the clicked
	// element (often a button this action removes from the DOM), and go-app drops
	// any Dispatch whose source element is no longer mounted — which would leave
	// the UI stuck on "resolving…" after a chooser.
	ctx.Async(func() {
		err := fn()
		g.dispatch(func(ctx app.Context) {
			g.busy = false
			if err != nil {
				g.status = err.Error()
			}
			g.afterAction()
			g.save(ctx)
		})
	})
}

// afterAction settles the phase after an engine mutation: the game may be won,
// the active house may have been cleared (a new turn began), or play continues.
func (g *game) afterAction() {
	g.clearSelection()
	g.settlePhase()
}

// settlePhase picks the resting interaction phase from the engine state: the game
// is over, a new turn needs a house, or play continues. Transient phases (picking
// a flank or a fight target) depend on a live selection, so they are never a
// resting phase and a resumed match always lands on one of these.
func (g *game) settlePhase() {
	switch {
	case g.g.Winner() >= 0:
		g.phase = phaseOver
	case g.g.State.ActiveHouse == engine.HouseNone:
		g.phase = phaseHouse
	default:
		g.phase = phaseMain
	}
}

// ---- selection ----
//
// Card click handlers take the card id as a parameter rather than capturing it in
// a per-card closure. go-app compares event handlers by function pointer, so
// every closure built from the same func literal looks identical to it and a
// captured id would never be refreshed when the board re-renders (e.g. when the
// turn flips). cardView instead keeps the id in a field and passes the live value
// to these methods on click.

// selectBoardID selects a card in play, deriving whether it belongs to the active
// player (actionable) or is another card (read-only).
func (g *game) selectBoardID(_ app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.phase == phaseFightTarget {
		return
	}
	g.sel, g.selKind, g.selHand, g.hasSel = id, g.boardKindOf(id), -1, true
	g.detailDef = nil
	g.status = ""
}

// selectHandID selects a card in the active player's hand, recovering its hand
// index from the id.
func (g *game) selectHandID(_ app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.phase == phaseFightTarget {
		return
	}
	idx := indexOfID(g.g.Hand(g.active()), id)
	if idx < 0 {
		return
	}
	g.sel, g.selKind, g.selHand, g.hasSel = id, selHand, idx, true
	g.detailDef = nil
	g.status = ""
}

// boardKindOf reports how a card in play should be treated when selected: one of
// the active player's own creatures/artifacts, or another (read-only) card.
func (g *game) boardKindOf(id engine.LocalID) selKind {
	active := g.active()
	switch {
	case containsID(g.g.Battleline(active), id):
		return selYourCreature
	case containsID(g.g.Artifacts(active), id):
		return selYourArtifact
	default:
		return selOther
	}
}

// ---- turn / play handlers ----

func (g *game) pickHouse(h engine.House) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if g.busy || g.choosing {
			return
		}
		p := g.active()
		g.runAction(ctx, func() error { return g.g.ChooseHouse(p, h) })
	}
}

func (g *game) endTurn(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing {
		return
	}
	p := g.active()
	g.runAction(ctx, func() error {
		g.g.EndTurn(p)
		g.g.BeginTurn(1 - p) // forge + start-of-turn triggers for the next player
		return nil
	})
}

// play resolves the selected hand card. Creatures ask for a flank first (unless
// the battleline is empty); everything else plays immediately.
func (g *game) play(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	p := g.active()
	idx := g.selHand
	def := g.g.Def(g.sel)
	switch def.Type {
	case engine.Creature:
		if len(g.g.Battleline(p)) == 0 {
			g.runAction(ctx, func() error { _, err := g.g.PlayCreature(p, idx, false); return err })
			return
		}
		g.phase = phaseFlank
	case engine.Artifact:
		g.runAction(ctx, func() error { _, err := g.g.PlayArtifact(p, idx); return err })
	case engine.Action:
		g.runAction(ctx, func() error { return g.g.PlayAction(p, idx) })
	case engine.Upgrade:
		g.runAction(ctx, func() error { _, err := g.g.PlayUpgrade(p, idx); return err })
	}
}

func (g *game) playFlank(left bool) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if g.busy || g.choosing || g.phase != phaseFlank || g.selKind != selHand {
			return
		}
		p, idx := g.active(), g.selHand
		g.runAction(ctx, func() error { _, err := g.g.PlayCreature(p, idx, left); return err })
	}
}

func (g *game) discard(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	p, idx := g.active(), g.selHand
	g.runAction(ctx, func() error { return g.g.DiscardFromHand(p, idx) })
}

// ---- drag and drop (hand → board) ----

// startHandDrag begins dragging a playable hand card. It selects the card so the
// detail panel and the drop share the same target, and marks a drag in progress
// so the board shows as a drop zone.
func (g *game) startHandDrag(ctx app.Context, id engine.LocalID) {
	if g.busy || g.choosing || g.choosingOption || g.phase != phaseMain {
		return
	}
	g.selectHandID(ctx, id)
	g.dragging = true
}

// endHandDrag clears the drag state when the pointer is released, whether or not
// the card landed on the board.
func (g *game) endHandDrag(_ app.Context, _ engine.LocalID) {
	g.dragging = false
}

// dropOnBoard plays the dragged hand card when it is released over the play area,
// following the normal flow — a creature still prompts for its flank. A card that
// cannot be played from hand is left where it is.
func (g *game) dropOnBoard(ctx app.Context, e app.Event) {
	e.PreventDefault()
	if !g.dragging {
		return
	}
	g.dragging = false
	if g.busy || g.choosing || g.choosingOption || g.phase != phaseMain || g.selKind != selHand {
		return
	}
	if !g.playableFromHand(g.sel) {
		return
	}
	g.play(ctx, e)
}

// ---- creature actions ----

func (g *game) reap(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selYourCreature {
		return
	}
	p, id := g.active(), g.sel
	g.runAction(ctx, func() error { return g.g.Reap(p, id) })
}

func (g *game) useAction(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain {
		return
	}
	if g.selKind != selYourCreature && g.selKind != selYourArtifact {
		return
	}
	p, id := g.active(), g.sel
	g.runAction(ctx, func() error { return g.g.UseAction(p, id) })
}

// startFight enters fight-target selection for the selected creature, after
// checking it can actually be used (so the player is not left picking a target
// for an exhausted or out-of-house attacker).
func (g *game) startFight(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.phase != phaseMain || g.selKind != selYourCreature {
		return
	}
	if err := g.g.CanUse(g.active(), g.sel); err != nil {
		g.status = err.Error()
		return
	}
	g.attacker = g.sel
	g.phase = phaseFightTarget
	g.status = ""
}

func (g *game) fightTargetID(ctx app.Context, defender engine.LocalID) {
	if g.busy || g.choosing || g.phase != phaseFightTarget {
		return
	}
	p, att := g.active(), g.attacker
	g.runAction(ctx, func() error { return g.g.Fight(p, att, defender) })
}

func (g *game) cancelTargeting(ctx app.Context, _ app.Event) {
	g.phase = phaseMain
	g.attacker = 0
}

// ---- engine chooser overlay ----

func (g *game) chooseCandidate(_ app.Context, id engine.LocalID) {
	if !g.choosing {
		return
	}
	select {
	case g.chooser.reply <- chooseReply{id: id, ok: true}:
	default:
	}
}

// chooseOptionIdx answers the current option prompt with option i. The index is
// stable per button position, so a captured value is safe here (unlike per-card
// closures).
func (g *game) chooseOptionIdx(i int) app.EventHandler {
	return func(_ app.Context, _ app.Event) {
		if !g.choosingOption {
			return
		}
		select {
		case g.chooser.optionReply <- i:
		default:
		}
	}
}

// onScorePillClick opens the out-of-play zone viewer for the clicked player. The
// player index is read from the pill's data attribute rather than captured in a
// closure, so the single stable handler stays valid across re-renders (go-app
// compares event handlers by function pointer). The viewer is read-only, so it
// stays available during a prompt — e.g. to inspect archives before deciding
// whether to take them into hand.
func (g *game) onScorePillClick(ctx app.Context, _ app.Event) {
	p, err := strconv.Atoi(ctx.JSSrc().Get("dataset").Get("player").String())
	if err != nil {
		return
	}
	g.zonesPlayer = p
}

// closeZones hides the out-of-play zone viewer.
func (g *game) closeZones(_ app.Context, _ app.Event) {
	g.zonesPlayer = -1
}

// stopClick keeps a click inside the zone panel from bubbling up to the
// backdrop's close handler, so only clicks outside the panel dismiss the viewer.
func (g *game) stopClick(_ app.Context, e app.Event) {
	e.Call("stopPropagation")
}

// ---- manual (sandbox) mode ----

// toggleManual turns the engine's sandbox mode on or off, lifting house
// restrictions and revealing the manual controls.
func (g *game) toggleManual(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing || g.choosingOption {
		return
	}
	g.g.SetManual(!g.g.Manual())
	g.save(ctx)
}

// manualMove moves the selected card to a resting zone, ignoring the normal rules.
func (g *game) manualMove(dest engine.ManualZone) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		if !g.hasSel || !g.g.Manual() {
			return
		}
		g.g.ManualMove(g.sel, dest)
		g.clearSelection()
		g.save(ctx)
	}
}

// manualReady clears the selected card's exhausted flag.
func (g *game) manualReady(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.g.ManualSetExhausted(g.sel, false)
	g.save(ctx)
}

// manualExhaust sets the selected card's exhausted flag.
func (g *game) manualExhaust(ctx app.Context, _ app.Event) {
	if !g.hasSel || !g.g.Manual() {
		return
	}
	g.g.ManualSetExhausted(g.sel, true)
	g.save(ctx)
}

// selectZoneCard selects a card shown in the zone viewer (in manual mode) and
// closes the viewer, so the manual controls can act on it (e.g. move it to hand).
func (g *game) selectZoneCard(_ app.Context, id engine.LocalID) {
	g.sel, g.selKind, g.selHand, g.hasSel = id, selOther, -1, true
	g.detailDef = nil
	g.zonesPlayer = -1
	g.status = ""
}

// openPicker opens the fuzzy card picker to add an arbitrary card to hand.
func (g *game) openPicker(_ app.Context, _ app.Event) {
	g.pickerOpen = true
	g.pickerQuery = ""
}

func (g *game) closePicker(_ app.Context, _ app.Event) { g.pickerOpen = false }

// pickerInput records the search box's text as the player types.
func (g *game) pickerInput(ctx app.Context, _ app.Event) {
	g.pickerQuery = ctx.JSSrc().Get("value").String()
}

// addPickedCard adds the chosen card definition to the active player's hand.
func (g *game) addPickedCard(def engine.CardDefinition) app.EventHandler {
	return func(ctx app.Context, _ app.Event) {
		g.g.ManualAddCard(def, g.active())
		g.pickerOpen = false
		g.save(ctx)
	}
}

// isInPlay reports whether a card is on either player's battleline or artifact row.
func (g *game) isInPlay(id engine.LocalID) bool {
	for p := 0; p < 2; p++ {
		if containsID(g.g.Battleline(p), id) || containsID(g.g.Artifacts(p), id) {
			return true
		}
	}
	return false
}

// onLogCardClick opens the detail panel on a card named in the log. The name is
// read from the clicked span's data attribute rather than captured in a closure,
// so the single stable handler stays valid as the log grows or resets (go-app
// compares event handlers by function pointer).
func (g *game) onLogCardClick(ctx app.Context, _ app.Event) {
	name := ctx.JSSrc().Get("dataset").Get("card").String()
	def, ok := g.defByName[name]
	if !ok {
		return
	}
	g.clearSelection()
	g.detailDef = def
}

func (g *game) restart(ctx app.Context, _ app.Event) {
	if g.busy || g.choosing {
		return
	}
	g.newMatch()
	g.save(ctx)
}
