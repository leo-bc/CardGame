package model

// ClientBattle :
type ClientBattle struct {
	Sides     []ClientBattleSide
	Cards     map[int]Card
	Turn      int
	IsStarted bool
}

// ClientBattleSide :
type ClientBattleSide struct {
	IsPlayer bool
	IsTurn   bool
	Cards    map[string][]CardSlot
}
