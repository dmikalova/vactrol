package web

import (
	"strings"

	"github.com/dmikalova/vactrol/internal/engine"
	"github.com/maxence-charriere/go-app/v11/pkg/app"
)

// cardView is a presentational component for a single card face. It carries no
// game logic: the parent supplies already-rendered strings, visual flags, and a
// click handler, so the same component renders a hand card, a creature in play,
// an artifact, or a targeting candidate.
type cardView struct {
	app.Compo

	ID         engine.LocalID
	Title      string
	HouseCls   string // house-derived border/background classes
	Stat       string // compact stat line (power, damage, Æmber…)
	Rules      string // rules/ability text for the face
	Kind       string // card type label shown at the foot
	Selected   bool
	Targetable bool
	Dimmed     bool
	// OnActivate is called with ID when the card is clicked; nil means the card is
	// not clickable. The id is passed rather than captured in the handler because
	// go-app compares event handlers by function pointer and would not refresh a
	// captured id when the board re-renders.
	OnActivate func(app.Context, engine.LocalID)
}

// onClick is a stable method (unlike a per-card closure) so go-app keeps it bound
// across re-renders; it reads the up-to-date ID and OnActivate fields at click
// time.
func (c *cardView) onClick(ctx app.Context, _ app.Event) {
	if c.OnActivate != nil {
		c.OnActivate(ctx, c.ID)
	}
}

func (c *cardView) Render() app.UI {
	clickable := c.OnActivate != nil
	cls := cx(
		"card",
		c.HouseCls,
		ifCls(c.Selected, "card--selected"),
		ifCls(c.Targetable, "card--targetable"),
		ifCls(c.Dimmed, "card--dimmed"),
		ifCls(clickable && !c.Targetable, "card--clickable"),
	)

	div := app.Div().Class(cls)
	if clickable {
		div = div.OnClick(c.onClick)
	}

	return div.Body(
		app.Div().Class("card-name").Text(c.Title),
		app.Div().Class("card-body").Body(
			app.If(c.Stat != "", func() app.UI {
				return app.Div().Class("card-stat").Text(c.Stat)
			}),
			app.If(c.Rules != "", func() app.UI {
				return app.Div().Class("card-rules").Text(c.Rules)
			}),
		),
		app.Div().Class("card-kind").Text(c.Kind),
	)
}

// cx joins non-empty class fragments with spaces.
func cx(parts ...string) string {
	kept := parts[:0]
	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}
	return strings.Join(kept, " ")
}

// ifCls returns cls when cond holds, otherwise the empty string.
func ifCls(cond bool, cls string) string {
	if cond {
		return cls
	}
	return ""
}
