package model

// State :
type State struct {
	Games   []Game
	Players []Player
	Cards   []Card
}

// Game :
type Game struct {
	Updateable
	Info GameInfo
}

// GameInfo :
type GameInfo struct {
	Title     string
	Players   []PlayerGameLink
	Battles   []Battle
	IsStarted bool
}

// Battle :
type Battle struct {
	Updateable
	Sides     []BattleSide
	Turn      int
	IsStarted bool
}

// BattleSide :
type BattleSide struct {
	PlayerID int
	Info     BattleSideInfo
}

// BattleSideInfo :
type BattleSideInfo struct {
	Cards map[string]*CardsContainer
}

// CardsContainer :
type CardsContainer struct {
	Slots []CardSlot
}

// Append :
func (container *CardsContainer) Append(slot CardSlot) {
	container.Slots = append(container.Slots, slot)
}

// Replace :
func (container *CardsContainer) Replace(slots []CardSlot) {
	container.Slots = slots
}

// CardSlot :
type CardSlot struct {
	CardID      int
	DamageTaken int
}

// PlayerGameLink :
type PlayerGameLink struct {
	PlayerID int
	CardIDs  []int
	IsReady  bool
}

// Player :
type Player struct {
	Name string
}

// Card :
type Card struct {
	Identity CardIdentityInfo
	Rank     CardRankInfo
	Attacks  []Attack
}

// CardIdentityInfo :
type CardIdentityInfo struct {
	Title       string
	Description string
	Type        string
	URL         string
}

// CardRankInfo :
type CardRankInfo struct {
	Ranking string
	Level   int
	HP      int
}

// Attack :
type Attack struct {
	Name   string
	Cost   int
	Damage int
}

// AttackInfo :
type AttackInfo struct {
	Source int
	Attack int
	Target int
}
