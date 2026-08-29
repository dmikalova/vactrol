package engine

import "testing"

func TestForgeKeyEffect(t *testing.T) {
	g := NewGame("A", "B", 1)
	ctx := &EffectContext{Resolver: g, Controller: 0}
	e := ForgeKey{}
	if e.Text() != "forge a key at current cost" {
		t.Errorf("text = %q", e.Text())
	}

	// Not enough Æmber: no key is forged.
	e.Resolve(ctx)
	if g.Keys(0) != 0 {
		t.Errorf("keys = %d, want 0 (could not afford)", g.Keys(0))
	}

	// Enough Æmber: one key is forged and its cost paid.
	g.State.Aember[0] = KeyCost + 2
	e.Resolve(ctx)
	if g.Keys(0) != 1 {
		t.Errorf("keys = %d, want 1", g.Keys(0))
	}
	if g.Aember(0) != 2 {
		t.Errorf("aember = %d, want 2 (paid the key cost)", g.Aember(0))
	}
}
