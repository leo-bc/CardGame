package model

// ClientBattle :
type ClientBattle struct {
	Sides []ClientBattleSide
	Cards map[int]Card
	Info  BattleInfo
}

// ClientBattleSide :
type ClientBattleSide struct {
	IsPlayer bool
	IsTurn   bool
	Cards    map[string][]CardSlot
}
