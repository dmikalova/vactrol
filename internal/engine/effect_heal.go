package engine

import "fmt"

// Heal removes damage from each creature its Target selects: a fixed Amount of
// damage, or all of it when Fully is set.
//
// Healing takes damage tokens off a creature. It can never remove more damage
// than is on the creature (a creature with no damage is unaffected), and it never
// changes a creature's power. Set exactly one of Amount or Fully — combining them
// is rejected when the card is built (see validate), never silently resolved.
type Heal struct {
	Amount int  // damage to remove; must be zero when Fully is set
	Fully  bool // remove all damage instead of a fixed Amount
	Target Target
}

// Text renders the effect, e.g. "heal 1 damage from each creature" or
// "fully heal each other friendly creature".
func (e Heal) Text() string {
	if e.Fully {
		return "fully heal " + e.Target.Text()
	}
	return fmt.Sprintf("heal %d damage from %s", e.Amount, e.Target.Text())
}

// Resolve removes the damage from each selected creature (all of it when Fully).
func (e Heal) Resolve(ctx *EffectContext) {
	for _, id := range e.Target.Select(ctx) {
		if e.Fully {
			ctx.Resolver.SetDamage(id, 0)
			continue
		}
		ctx.Resolver.SetDamage(id, ctx.Resolver.Damage(id)-e.Amount)
	}
}

// validate rejects a Heal that sets both a fixed Amount and Fully, since the two
// are different ways to say how much to heal and combining them is ambiguous.
func (e Heal) validate() error {
	if e.Fully && e.Amount != 0 {
		return fmt.Errorf("heal: set Amount or Fully, not both (got Amount=%d, Fully=true)", e.Amount)
	}
	return nil
}
