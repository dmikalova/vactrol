package engine

import "testing"

// TestSetAemberClampsAtZero exercises the Resolver's Æmber setter, including its
// floor at zero (a pool can never go negative).
func TestSetAemberClampsAtZero(t *testing.T) {
	g := NewGame("A", "B", 1)
	g.SetAember(0, 5)
	if g.Aember(0) != 5 {
		t.Errorf("Aember = %d, want 5", g.Aember(0))
	}
	g.SetAember(0, -3) // negative clamps to zero
	if g.Aember(0) != 0 {
		t.Errorf("Aember after negative set = %d, want 0", g.Aember(0))
	}
}
