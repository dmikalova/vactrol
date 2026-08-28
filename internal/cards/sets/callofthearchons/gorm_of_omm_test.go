package callofthearchons

import (
	"testing"

	"github.com/dmikalova/vactrol/internal/cards/cardtest"
	"github.com/dmikalova/vactrol/internal/engine"
)

// Gorm of Omm sacrifices itself and destroys another artifact.
func TestGormOfOmm(t *testing.T) {
	g := cardtest.Started(t, engine.Sanctum)
	gorm := g.AddArtifact(GormOfOmm, 0)
	g.AddArtifact(Cannon, 1) // the only other artifact, so it is the forced choice

	if err := g.UseAction(0, gorm); err != nil {
		t.Fatalf("UseAction: %v", err)
	}
	if len(g.Artifacts(0)) != 0 {
		t.Error("Gorm of Omm should sacrifice itself")
	}
	if len(g.Artifacts(1)) != 0 {
		t.Errorf("Gorm of Omm should destroy the other artifact; remaining %v", g.Artifacts(1))
	}
}
