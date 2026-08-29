package web

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// Render draws the whole client. It runs on both the server (prerender) and the
// client; before OnMount seeds the match on the client, g is nil, so a lightweight
// placeholder is shown.
func (g *game) Render() app.UI {
	if g.g == nil {
		return app.Div().Class("loading").Text("Loading Vactrol…")
	}

	return app.Div().Class("app").Body(
		app.Div().Class("board-area").Body(g.boardArea()...),
		app.Div().Class("sidebar").Body(
			g.brandBar(),
			g.logPanel(),
			g.detailPanel(),
			g.controls(),
		),
		app.If(g.zonesPlayer >= 0, func() app.UI { return g.zonesOverlay() }),
		app.If(g.phase == phaseOver, func() app.UI { return g.overBanner() }),
	)
}

// brandBar is the slim top of the sidebar: the title, a transient status, and the
// New game button.
func (g *game) brandBar() app.UI {
	return app.Div().Class("brandbar").Body(
		app.Span().Class("brand-title").Text("Vactrol"),
		app.If(g.busy && !g.choosing && !g.choosingOption, func() app.UI {
			return app.Span().Class("badge-busy").Text("resolving…")
		}),
		app.If(g.status != "", func() app.UI {
			return app.Span().Class("badge-error").Text(g.status)
		}),
		app.Div().Class("spacer"),
		app.Button().Class("btn-nav").Text("New game").
			Disabled(g.busy || g.choosing || g.choosingOption).
			OnClick(g.restart),
	)
}

// ---- board ----

// boardArea renders both battlelines facing each other (opponent on top), each
// player's score, and the active player's hand. Card rows squeeze to the height
// available and scroll horizontally when full. The play rows between the two
// score pills form the drop zone for playing a card dragged from hand.
func (g *game) boardArea() []app.UI {
	p := g.active()
	opp := 1 - p
	playZone := app.Div().
		Class(cx("play-zone", ifCls(g.dragging, "play-zone--drop"))).
		OnDrop(g.dropOnBoard).
		Body(
			g.renderRow(g.g.PlayerName(opp)+"'s artifacts", g.g.Artifacts(opp), selOther),
			g.renderRow(g.g.PlayerName(opp)+"'s creatures", g.g.Battleline(opp), selOther),
			app.Div().Class("midline"),
			g.renderRow(g.g.PlayerName(p)+"'s creatures", g.g.Battleline(p), selYourCreature),
			g.renderRow(g.g.PlayerName(p)+"'s artifacts", g.g.Artifacts(p), selYourArtifact),
		)
	return []app.UI{
		g.scorePill(opp),
		playZone,
		g.scorePill(p),
		g.renderHand(),
	}
}

func (g *game) scorePill(player int) app.UI {
	active := player == g.active()
	parts := []string{
		fmt.Sprintf("%d Æmber", g.g.Aember(player)),
		fmt.Sprintf("%d/%d keys", g.g.Keys(player), engine.KeysToWin),
	}
	detail := []app.UI{app.Text(strings.Join(parts, " • "))}
	if active && g.g.State.ActiveHouse != engine.HouseNone {
		h := g.g.State.ActiveHouse
		detail = append(detail,
			app.Text(" • house "),
			app.Span().Class(cx("score-house", houseAccent(h))).Text(h.String()),
		)
	}
	zones := fmt.Sprintf("deck %d • discard %d • archives %d • purge %d",
		len(g.g.Deck(player)), len(g.g.Discard(player)),
		len(g.g.Archives(player)), len(g.g.Purge(player)))
	cls := cx("score-pill", ifCls(active, "score-pill-active"), ifCls(!active, "score-pill-idle"))
	return app.Div().Class(cls).
		DataSet("player", strconv.Itoa(player)).
		OnClick(g.onScorePillClick).
		Body(
			app.Span().Class("score-name").Text(g.g.PlayerName(player)),
			app.Span().Class("score-detail").Body(detail...),
			app.Span().Class("score-zones").Text(zones),
		)
}

func (g *game) renderRow(label string, ids []engine.LocalID, boardKind selKind) app.UI {
	// The row always reserves card height (board-row--fill) so playing or losing a
	// card does not resize the row and shift the rest of the board.
	return app.Div().Class("board-row board-row--fill").Body(
		app.Div().Class("row-label").Text(label),
		app.If(len(ids) == 0, func() app.UI {
			return app.Div().Class("row-empty").Text("—")
		}).Else(func() app.UI {
			return app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI { return g.renderCard(ids[i], boardKind) }),
			)
		}),
	)
}

func (g *game) renderCard(id engine.LocalID, boardKind selKind) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, boardKind)
	return &cardView{
		ID:         id,
		Title:      def.Name,
		HouseCls:   houseClasses(def.House),
		Stat:       g.statLine(id),
		Rules:      g.faceRules(id),
		Kind:       string(def.Type),
		Selected:   g.hasSel && g.sel == id,
		Targetable: targetable,
		Dimmed:     dimmed,
		OnActivate: activate,
	}
}

// cardVisual decides how a card (in hand or in play) responds and looks in the
// current mode. It returns the click handler (nil when the card is not
// clickable), whether the card is a highlighted action target, and whether it is
// lowlighted (dimmed) as an invalid choice. During a chooser or fight-target
// prompt only the eligible cards are highlighted and the rest dimmed; in ordinary
// play, cards the active player cannot act with (wrong house, exhausted, or
// unplayable from hand) are dimmed so the usable ones stand out.
func (g *game) cardVisual(id engine.LocalID, kind selKind) (activate func(app.Context, engine.LocalID), targetable, dimmed bool) {
	switch {
	case g.choosing:
		// A chooser runs on a background goroutine, so g.busy is also set; the
		// choosing case must come first or the candidates would not be clickable.
		if containsID(g.chooserCandidates, id) {
			return g.chooseCandidate, true, false
		}
		return nil, false, true
	case g.busy:
		return nil, false, false
	case g.phase == phaseFightTarget:
		if kind == selOther && g.isEnemyCreature(id) {
			return g.fightTargetID, true, false
		}
		return nil, false, true
	case kind == selHand:
		return g.selectHandID, false, !g.playableFromHand(g.g.Def(id))
	case kind == selYourCreature, kind == selYourArtifact:
		return g.selectBoardID, false, !g.actionable(id, kind)
	default:
		return g.selectBoardID, false, false
	}
}

// actionable reports whether the active player can act with one of their own
// cards this turn, so the ones they cannot use are lowlighted. Creatures defer to
// the engine's CanUse (house, readiness, ownership); an artifact needs an action
// ability, readiness, and the active house.
func (g *game) actionable(id engine.LocalID, kind selKind) bool {
	switch kind {
	case selYourCreature:
		return g.g.CanUse(g.active(), id) == nil
	case selYourArtifact:
		def := g.g.Def(id)
		inHouse := g.g.State.ActiveHouse == engine.HouseNone || def.House == g.g.State.ActiveHouse
		return inHouse && g.g.HasAction(id) && !g.g.Exhausted(id)
	default:
		return true
	}
}

// ---- hand ----

func (g *game) renderHand() app.UI {
	p := g.active()
	ids := g.sortedHand(p)
	cls := "board-row board-row--fill"
	return app.Div().Class(cls).Body(
		app.Div().Class("row-label").Text(g.g.PlayerName(p)+"'s hand"),
		app.Div().Class("card-strip").Body(
			app.Range(ids).Slice(func(i int) app.UI { return g.renderHandCard(ids[i]) }),
		),
	)
}

// sortedHand returns the player's hand ids ordered by house, then card type, then
// name, so the hand reads consistently. It sorts a copy — the engine's own hand
// order (which play/discard index into) is untouched, so selection still maps to
// the right card.
func (g *game) sortedHand(p int) []engine.LocalID {
	return g.sortByHouseTypeName(g.g.Hand(p))
}

// sortByHouseTypeName returns a copy of ids ordered by house, then card type,
// then name — the stable reading order shared by the hand and the deck view. The
// deck in particular must not reveal its shuffled order, so it is always sorted.
func (g *game) sortByHouseTypeName(ids []engine.LocalID) []engine.LocalID {
	ids = append([]engine.LocalID(nil), ids...)
	sort.SliceStable(ids, func(i, j int) bool {
		a, b := g.g.Def(ids[i]), g.g.Def(ids[j])
		if a.House != b.House {
			return a.House < b.House
		}
		if a.Type != b.Type {
			return a.Type < b.Type
		}
		return a.Name < b.Name
	})
	return ids
}

func (g *game) renderHandCard(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	activate, targetable, dimmed := g.cardVisual(id, selHand)
	draggable := !g.busy && !g.choosing && !g.choosingOption &&
		g.phase == phaseMain && g.playableFromHand(def)
	return &cardView{
		ID:          id,
		Title:       def.Name,
		HouseCls:    houseClasses(def.House),
		Stat:        handStat(def),
		Rules:       engine.RenderCardRules(def),
		Kind:        string(def.Type),
		Selected:    g.hasSel && g.sel == id,
		Targetable:  targetable,
		Dimmed:      dimmed,
		OnActivate:  activate,
		Draggable:   draggable,
		OnDragStart: g.startHandDrag,
		OnDragEnd:   g.endHandDrag,
	}
}

// ---- sidebar controls ----

// controls is the bottom of the sidebar: the contextual controls (house picker or
// action bar) plus End turn. House selection has no End turn — a house must be
// chosen first.
func (g *game) controls() app.UI {
	// A labeled option prompt (e.g. "take archives?") shows its choices as buttons.
	if g.choosingOption {
		return app.Div().Class("controls").Body(g.optionChooser())
	}
	// While an engine chooser waits, the controls become the prompt itself: a
	// green call to action to click one of the highlighted cards. The choice is
	// mandatory, so there is no cancel.
	if g.choosing {
		return app.Div().Class("controls").Body(
			app.Div().Class("prompt").Text(g.chooserPrompt),
		)
	}
	if g.phase == phaseHouse {
		return app.Div().Class("controls").Body(g.housePicker())
	}
	// Mid-action selections (choosing a flank or a fight target) show only their
	// own controls — End turn is withheld until the action finishes or is cancelled.
	if g.phase == phaseFlank || g.phase == phaseFightTarget {
		return app.Div().Class("controls").Body(g.actionBar())
	}
	return app.Div().Class("controls").Body(
		g.actionBar(),
		btn("End turn", g.endTurn, "btn-secondary"),
	)
}

func (g *game) housePicker() app.UI {
	houses := g.deckHouses[g.active()]
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text("Choose a house"),
		app.Range(houses).Slice(func(i int) app.UI {
			h := houses[i]
			return app.Button().
				Class(cx("house-btn", houseAccent(h))).
				Text(h.String()).
				OnClick(g.pickHouse(h))
		}),
	)
}

// optionChooser renders a labeled multiple-choice prompt. When every option is a
// house (a "choose a house" effect such as Blinding Light) it shows colored house
// buttons like the start-of-turn picker; otherwise plain primary buttons.
func (g *game) optionChooser() app.UI {
	if g.houseOptions() {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Choose a house:"),
			app.Range(g.optionLabels).Slice(func(i int) app.UI {
				h, _ := engine.ParseHouse(g.optionLabels[i])
				return app.Button().
					Class(cx("house-btn", houseAccent(h))).
					Text(g.optionLabels[i]).
					OnClick(g.chooseOptionIdx(i))
			}),
		)
	}
	return app.Div().Class("btn-col").Body(
		app.Div().Class("prompt").Text(g.optionPrompt),
		app.Range(g.optionLabels).Slice(func(i int) app.UI {
			return btn(g.optionLabels[i], g.chooseOptionIdx(i), "btn-primary")
		}),
	)
}

// houseOptions reports whether every current option label names a house, so the
// prompt can be shown as the colored house picker.
func (g *game) houseOptions() bool {
	if len(g.optionLabels) == 0 {
		return false
	}
	for _, label := range g.optionLabels {
		if _, ok := engine.ParseHouse(label); !ok {
			return false
		}
	}
	return true
}

func (g *game) actionBar() app.UI {
	switch {
	case g.phase == phaseFlank:
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text("Play "+g.g.Def(g.sel).Name+" to which flank?"),
			btn("Left flank", g.playFlank(true), "btn-primary"),
			btn("Right flank", g.playFlank(false), "btn-primary"),
			btn("Cancel", g.cancelTargeting, "btn-secondary"),
		)
	case g.phase == phaseFightTarget:
		return app.Div().Class("btn-col").Body(
			app.Div().Class("prompt").Text("Pick an enemy creature to fight"),
			btn("Cancel", g.cancelTargeting, "btn-secondary"),
		)
	case !g.hasSel:
		return app.Div().Class("hint").Text("Select a card to inspect it and choose an action.")
	}

	switch g.selKind {
	case selHand:
		return g.handActions()
	case selYourCreature:
		return g.creatureActions()
	case selYourArtifact:
		return g.artifactActions()
	default:
		return app.Div().Class("hint").Text("Read-only — this is your opponent's card.")
	}
}

func (g *game) handActions() app.UI {
	def := g.g.Def(g.sel)
	if h := g.g.State.ActiveHouse; h != engine.HouseNone && def.House != h {
		return app.Div().Class("hint").Text("Not the active house — you cannot play this from hand.")
	}
	if def.Type == engine.Upgrade && !g.anyCreatureInPlay() {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text(def.Name),
			app.Div().Class("hint").Text("No creature in play to attach this upgrade to."),
			btn("Discard", g.discard, "btn-danger"),
		)
	}
	return app.Div().Class("btn-col").Body(
		app.Div().Class("section-title").Text(def.Name),
		btn("Play", g.play, "btn-primary"),
		btn("Discard", g.discard, "btn-danger"),
	)
}

func (g *game) creatureActions() app.UI {
	if err := g.g.CanUse(g.active(), g.sel); err != nil {
		return app.Div().Class("btn-col").Body(
			app.Div().Class("section-title").Text(g.g.Def(g.sel).Name),
			app.Div().Class("hint").Text("Cannot act: "+err.Error()+"."),
		)
	}
	items := []app.UI{
		app.Div().Class("section-title").Text(g.g.Def(g.sel).Name),
		btn("Reap", g.reap, "btn-primary"),
	}
	// Fight is offered only when a legal target exists (e.g. with no enemy
	// creatures, a ready Valdr can still reap but has nothing to fight).
	if len(g.g.FightTargets(g.active(), g.sel)) > 0 {
		items = append(items, btn("Fight", g.startFight, "btn-primary"))
	}
	if g.g.HasAction(g.sel) {
		items = append(items, btn("Action", g.useAction, "btn-primary"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) artifactActions() app.UI {
	name := g.g.Def(g.sel).Name
	if !g.g.HasAction(g.sel) {
		return app.Div().Class("hint").Text(name + " has no action ability.")
	}
	items := []app.UI{app.Div().Class("section-title").Text(name)}
	if g.g.Exhausted(g.sel) {
		items = append(items, app.Div().Class("hint").Text("Exhausted."))
	} else {
		items = append(items, btn("Action", g.useAction, "btn-primary"))
	}
	return app.Div().Class("btn-col").Body(items...)
}

func (g *game) detailPanel() app.UI {
	body := []app.UI{app.Div().Class("panel-label").Text("Card detail")}
	switch {
	case g.detailDef != nil:
		// A card referenced from the log: its printed text, without live state.
		body = append(body, app.Div().Class("detail-text").Text(engine.RenderCardText(g.detailDef)))
	case g.hasSel:
		body = append(body, app.Div().Class("detail-text").Text(engine.RenderCardText(g.g.Def(g.sel))))
		if a := g.g.AmberOn(g.sel); a > 0 {
			body = append(body, app.Div().Class("detail-amber").Text(fmt.Sprintf("Captured Æmber: %d", a)))
		}
		// Show attached upgrades so their effect is visible from the creature they
		// are on (upgrades are not separately clickable).
		for _, up := range g.g.Upgrades(g.sel) {
			updef := g.g.Def(up)
			body = append(body,
				app.Div().Class("detail-sub").Text("↳ "+updef.Name),
				app.Div().Class("detail-text").Text(engine.RenderCardText(updef)),
			)
		}
	default:
		body = append(body, app.Div().Class("detail-empty").Text("No card selected."))
	}
	return app.Div().Class("details").Body(body...)
}

func (g *game) logPanel() app.UI {
	return app.Div().Class("log").Body(
		app.Div().Class("panel-label").Text("Log"),
		app.Div().Class("log-list").ID("gamelog").Body(
			app.Range(g.g.Log).Slice(func(i int) app.UI {
				return app.Div().Body(g.logSegments(g.g.Log[i])...)
			}),
		),
	)
}

// logSegments splits a log line into plain text and clickable spans for the card
// names it mentions, so a mentioned card can be opened in the detail panel
// without changing the engine's log strings.
func (g *game) logSegments(line string) []app.UI {
	var out []app.UI
	var plain strings.Builder
	flush := func() {
		if plain.Len() > 0 {
			out = append(out, app.Text(plain.String()))
			plain.Reset()
		}
	}
	for i := 0; i < len(line); {
		if name := g.cardNameAt(line, i); name != "" {
			flush()
			out = append(out, app.Span().
				Class("log-card").
				DataSet("card", name).
				OnClick(g.onLogCardClick).
				Text(name))
			i += len(name)
			continue
		}
		plain.WriteByte(line[i])
		i++
	}
	flush()
	return out
}

// cardNameAt returns the longest known card name that begins at line[i] on word
// boundaries, or "" if none.
func (g *game) cardNameAt(line string, i int) string {
	if i > 0 && isWordByte(line[i-1]) {
		return "" // in the middle of a word — not a name boundary
	}
	best := ""
	for name := range g.defByName {
		if len(name) <= len(best) || !strings.HasPrefix(line[i:], name) {
			continue
		}
		if end := i + len(name); end < len(line) && isWordByte(line[end]) {
			continue // the name is only a prefix of a longer word
		}
		best = name
	}
	return best
}

func isWordByte(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b >= '0' && b <= '9'
}

// ---- overlays ----

func (g *game) overBanner() app.UI {
	winner := g.g.Winner()
	return app.Div().Class("over-backdrop").Body(
		app.Div().Class("over-card").Body(
			app.Div().Class("over-title").Text(g.g.PlayerName(winner)+" wins!"),
			app.Button().Class("btn-primary").Text("New game").OnClick(g.restart),
		),
	)
}

// zonesOverlay shows a single player's zones — their deck, discard, archives,
// and purge piles — as read-only card strips, so cards outside the board can be
// inspected. It is a modal dismissed by the ✕ in the corner or by clicking the
// backdrop outside the panel.
func (g *game) zonesOverlay() app.UI {
	p := g.zonesPlayer
	return app.Div().Class("over-backdrop").OnClick(g.closeZones).Body(
		app.Div().Class("zones-panel").OnClick(g.stopClick).Body(
			app.Button().Class("zones-close").Text("✕").OnClick(g.closeZones),
			app.Div().Class("over-title").Text(g.g.PlayerName(p)+"'s Zones"),
			g.zoneRow("Deck", g.sortByHouseTypeName(g.g.Deck(p))),
			g.zoneRow("Discard", g.g.Discard(p)),
			g.zoneRow("Archives", g.g.Archives(p)),
			g.zoneRow("Purge", g.g.Purge(p)),
		),
	)
}

func (g *game) zoneRow(label string, ids []engine.LocalID) app.UI {
	return app.Div().Class("zone-row").Body(
		app.Div().Class("row-label").Text(fmt.Sprintf("%s (%d)", label, len(ids))),
		app.If(len(ids) == 0, func() app.UI {
			return app.Div().Class("row-empty").Text("—")
		}).Else(func() app.UI {
			return app.Div().Class("card-strip").Body(
				app.Range(ids).Slice(func(i int) app.UI { return g.renderZoneCard(ids[i]) }),
			)
		}),
	)
}

// renderZoneCard renders a read-only card face for a card in an out-of-play zone.
func (g *game) renderZoneCard(id engine.LocalID) app.UI {
	def := g.g.Def(id)
	return &cardView{
		ID:       id,
		Title:    def.Name,
		HouseCls: houseClasses(def.House),
		Stat:     handStat(def),
		Rules:    engine.RenderCardRules(def),
		Kind:     string(def.Type),
	}
}

// ---- small helpers ----

func btn(label string, h app.EventHandler, class string) app.UI {
	return app.Button().Class(class).Text(label).OnClick(h)
}

func (g *game) statLine(id engine.LocalID) string {
	def := g.g.Def(id)
	var parts []string
	if def.Type == engine.Creature {
		parts = append(parts, fmt.Sprintf("%d power", g.g.Power(id)))
		if d := g.g.Damage(id); d > 0 {
			parts = append(parts, fmt.Sprintf("%d dmg", d))
		}
		if a := g.g.Armor(id); a > 0 {
			parts = append(parts, fmt.Sprintf("%d armor", a))
		}
	}
	if a := g.g.AmberOn(id); a > 0 {
		parts = append(parts, fmt.Sprintf("%d Æ", a))
	}
	if g.g.Exhausted(id) {
		parts = append(parts, "exhausted")
	}
	return strings.Join(parts, " • ")
}

func handStat(def *engine.CardDefinition) string {
	var parts []string
	if def.Type == engine.Creature {
		parts = append(parts, fmt.Sprintf("%d power", def.Power))
		if def.Armor > 0 {
			parts = append(parts, fmt.Sprintf("%d armor", def.Armor))
		}
	}
	if def.AemberBonus > 0 {
		parts = append(parts, fmt.Sprintf("+%d Æ", def.AemberBonus))
	}
	return strings.Join(parts, " • ")
}

func (g *game) faceRules(id engine.LocalID) string {
	var lines []string
	if s := engine.RenderCardRules(g.g.Def(id)); s != "" {
		lines = append(lines, s)
	}
	for _, up := range g.g.Upgrades(id) {
		lines = append(lines, "↳ "+g.g.Def(up).Name)
	}
	return strings.Join(lines, "\n")
}

func (g *game) playableFromHand(def *engine.CardDefinition) bool {
	h := g.g.State.ActiveHouse
	if h != engine.HouseNone && def.House != h {
		return false
	}
	if def.Type == engine.Upgrade && !g.anyCreatureInPlay() {
		return false
	}
	return true
}

// anyCreatureInPlay reports whether either battleline holds a creature. An
// upgrade attaches to a creature (on either side), so it cannot be played into an
// empty board.
func (g *game) anyCreatureInPlay() bool {
	return len(g.g.Battleline(0)) > 0 || len(g.g.Battleline(1)) > 0
}

func (g *game) isEnemyCreature(id engine.LocalID) bool {
	if g.g.Def(id).Type != engine.Creature {
		return false
	}
	for _, c := range g.g.Battleline(1 - g.active()) {
		if c == id {
			return true
		}
	}
	return false
}

func containsID(ids []engine.LocalID, id engine.LocalID) bool {
	for _, x := range ids {
		if x == id {
			return true
		}
	}
	return false
}

func indexOfID(ids []engine.LocalID, id engine.LocalID) int {
	for i, x := range ids {
		if x == id {
			return i
		}
	}
	return -1
}
