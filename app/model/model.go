package model

// State :
type State struct {
	Games   []Game
	Players []Player
	Cards   []Card
}

// IDable :
type IDable struct {
	ID int
}

// GetID :
func (id IDable) GetID() int {
	return id.ID
}

// Game :
type Game struct {
	IDable
	Info GameInfo
}

// GameInfo :
type GameInfo struct {
	Title   string
	Players []PlayerGameLink
}

// PlayerGameLink :
type PlayerGameLink struct {
	PlayerID int
	CardIDs  []int
}

// Player :
type Player struct {
	IDable
	Info PlayerInfo
}

// PlayerInfo :
type PlayerInfo struct {
	Name string
}

// Card :
type Card struct {
	IDable
	Info CardInfo
}

// CardInfo :
type CardInfo struct {
	Title string
	HP    int
}
