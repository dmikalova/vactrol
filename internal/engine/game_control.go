package engine

// takeControl moves a creature into controller's battleline without changing its
// owner. Ownership is immutable in KeyForge and still decides the out-of-play zone
// the card returns to; control is the battleline/current-player relationship.
func (g *Game) takeControl(id LocalID, controller int) {
	if g.inPlay(id) && g.cat.def(id).Type == Creature {
		g.removeFromPlay(id)
		controlPlus := uint8(0)
		if controller != g.owner(id) {
			controlPlus = uint8(controller + 1)
		}
		g.State.Cards[id].ControlPlus = controlPlus
		g.State.Battleline[controller].add(id)
		g.logf("%s takes control of %s", g.names[controller], g.Name(id))
	}
}

// revertControlFromUpgrade restores a controlled host to its owner's battleline
// when the Upgrade that took control of it leaves play. If the host has already
// left play, its resetCore call will clear the override instead.
func (g *Game) revertControlFromUpgrade(host, upgrade LocalID) {
	if g.cat.def(upgrade).Static.TakesControl && g.inPlay(host) &&
		g.State.Cards[host].ControlPlus != 0 && g.controller(host) == g.owner(upgrade) {
		owner := g.owner(host)
		g.removeFromPlay(host)
		g.State.Cards[host].ControlPlus = 0
		g.State.Battleline[owner].add(host)
		g.logf("%s returns to %s's control", g.Name(host), g.names[owner])
	}
}
