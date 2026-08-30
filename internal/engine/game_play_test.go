package engine

import "testing"

func TestCanPlay(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	creat := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, creat); err != nil {
		t.Errorf("playable creature = %v, want nil", err)
	}
	if err := g.CanPlay(1, creat); err != ErrNotActivePlayer {
		t.Errorf("wrong player = %v, want ErrNotActivePlayer", err)
	}

	off := g.AddToHand(NewCard("off", Sanctum, Creature, Common, WithPower(1)), 0)
	if err := g.CanPlay(0, off); err != ErrWrongHouse {
		t.Errorf("off-house = %v, want ErrWrongHouse", err)
	}

	up := g.AddToHand(NewCard("up", Brobnar, Upgrade, Common), 0)
	if err := g.CanPlay(0, up); err != ErrNoTarget {
		t.Errorf("hostless upgrade = %v, want ErrNoTarget", err)
	}

	g.State.Winner = 1
	if err := g.CanPlay(0, creat); err != ErrGameOver {
		t.Errorf("game over = %v, want ErrGameOver", err)
	}
}

func TestCanPlayRestrictions(t *testing.T) {
	// Creatures barred: a card that stops the player playing creatures.
	g := started(t)
	g.AddToBattleline(NewCard("Blocker", Brobnar, Creature, Common, WithPower(1),
		WithRestrictions(Restrictions{CannotPlay: Creature})), 0)
	c := g.AddToHand(testCreature("c", 3), 0)
	if err := g.CanPlay(0, c); err != ErrCannotPlayCreature {
		t.Errorf("creatures barred = %v, want ErrCannotPlayCreature", err)
	}

	// Card-play limit reached this turn.
	g2 := started(t)
	g2.AddToBattleline(NewCard("Imp", Brobnar, Creature, Common, WithPower(1),
		WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: Controller, Amount: 2}})), 0)
	g2.State.CardsPlayedThisTurn[0] = 2
	c2 := g2.AddToHand(testCreature("c2", 3), 0)
	if err := g2.CanPlay(0, c2); err != ErrCardPlayLimit {
		t.Errorf("limit reached = %v, want ErrCardPlayLimit", err)
	}
}

func TestCardPlayLimit(t *testing.T) {
	g := started(t) // player 0 active, Brobnar
	// Player 1 controls an Ember-Imp-like card limiting player 0 to two plays.
	g.AddToBattleline(NewCard("imp", Dis, Creature, Common, WithPower(2),
		WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: Opponent, Amount: 2}})), 1)
	g.AddToHand(testCreature("c0", 3), 0)
	g.AddToHand(testCreature("c1", 3), 0)
	g.AddToHand(testCreature("c2", 3), 0)
	g.AddToHand(exBruteStrength(), 0) // an upgrade

	if _, err := g.PlayCreature(0, handIdx(g, 0, "c0"), false); err != nil {
		t.Fatalf("first play: %v", err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c1"), false); err != nil {
		t.Fatalf("second play: %v", err)
	}
	// The third card play is barred, whatever its type.
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != ErrCardPlayLimit {
		t.Errorf("third creature play = %v, want ErrCardPlayLimit", err)
	}
	if _, err := g.PlayUpgrade(0, handIdx(g, 0, "Brute Strength")); err != ErrCardPlayLimit {
		t.Errorf("upgrade play = %v, want ErrCardPlayLimit", err)
	}

	// A new turn resets the count.
	g.EndTurn(0)
	g.BeginTurn(0)
	if err := g.ChooseHouse(0, Brobnar); err != nil {
		t.Fatal(err)
	}
	if _, err := g.PlayCreature(0, handIdx(g, 0, "c2"), false); err != nil {
		t.Errorf("play after new turn = %v, want nil (count reset)", err)
	}
}

func TestPlayCardLimitTargets(t *testing.T) {
	cases := []struct {
		name         string
		player       Player
		limited      int
		unrestricted int
	}{
		{name: "controller", player: Controller, limited: 0, unrestricted: 1},
		{name: "opponent", player: Opponent, limited: 1, unrestricted: 0},
		{name: "each player", player: EachPlayer, limited: 0, unrestricted: -1},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			g := started(t)
			g.AddToBattleline(NewCard("limit", Brobnar, Creature, Common, WithPower(1),
				WithRestrictions(Restrictions{PlayCardLimit: PlayCardLimit{Player: tc.player, Amount: 2}})), 0)
			g.State.CardsPlayedThisTurn[tc.limited] = 2
			if !g.cannotPlayCard(tc.limited) {
				t.Errorf("player %d should be limited", tc.limited)
			}
			if tc.unrestricted >= 0 && g.cannotPlayCard(tc.unrestricted) {
				t.Errorf("player %d should not be limited", tc.unrestricted)
			}
		})
	}
	if (PlayCardLimit{}).affects(0, 0) {
		t.Error("an unset play-card limit should affect no player")
	}
}
