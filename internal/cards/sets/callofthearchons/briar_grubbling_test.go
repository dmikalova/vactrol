package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Briar Grubbling has Hazardous 5: a creature that attacks it is dealt 5 damage
// before fight damage, enough to destroy most attackers before combat.
func TestBriarGrubbling(t *testing.T) {
	g := cardtest.Started(t, engine.Mars)
	grub := g.AddToBattleline(BriarGrubbling, 1) // 2 power, Hazardous 5
	attacker := g.AddToBattleline(cardtest.Vanilla("Attacker", engine.Mars, 4), 0)

	if err := g.Fight(0, attacker, grub); err != nil {
		t.Fatalf("Fight: %v", err)
	}
	// The 4-power attacker is destroyed by Hazardous 5 before combat, so the
	// grubbling exchanges no fight damage and survives unharmed.
	if len(g.Battleline(0)) != 0 {
		t.Error("attacker should be destroyed by Hazardous 5")
	}
	if g.Damage(grub) != 0 {
		t.Errorf("grubbling damage = %d, want 0 (no combat occurred)", g.Damage(grub))
	}
}
