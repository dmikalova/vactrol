package engine

import "testing"

func TestHouseString(t *testing.T) {
	if Brobnar.String() != "Brobnar" {
		t.Errorf("Brobnar.String() = %q", Brobnar.String())
	}
	if House(99).String() != "Unknown" {
		t.Errorf("House(99).String() = %q, want Unknown", House(99).String())
	}
}

func TestParseHouse(t *testing.T) {
	// Case-insensitive match.
	if h, ok := ParseHouse("bRoBnAr"); !ok || h != Brobnar {
		t.Errorf("ParseHouse(bRoBnAr) = %v, %v; want Brobnar, true", h, ok)
	}
	// "None" is not a selectable house.
	if h, ok := ParseHouse("none"); ok || h != HouseNone {
		t.Errorf("ParseHouse(none) = %v, %v; want HouseNone, false", h, ok)
	}
	// Unknown name.
	if h, ok := ParseHouse("nonsense"); ok || h != HouseNone {
		t.Errorf("ParseHouse(nonsense) = %v, %v; want HouseNone, false", h, ok)
	}
}
