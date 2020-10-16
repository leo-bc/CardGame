package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

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

// Updateable :
type Updateable struct {
	IsUpdated bool
}

// SetUpdated :
func (updateable *Updateable) SetUpdated() {
	updateable.IsUpdated = true
}

// TriggerUpdate :
func (updateable *Updateable) TriggerUpdate() bool {
	if updateable.IsUpdated {
		updateable.IsUpdated = false
		return true
	}
	return false
}

// Game :
type Game struct {
	IDable
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
	Bench     []CardSlot
	Hand      []CardSlot
	TakePile  []CardSlot
	ThrowPile []CardSlot
}

// PlayerBattleSideInfo :
type PlayerBattleSideInfo struct {
	IsPlayer bool
	IsTurn   bool
	Info     BattleSideInfo
}

// CardSlot :
type CardSlot struct {
	CardID int
}

// PlayerGameLink :
type PlayerGameLink struct {
	PlayerID int
	CardIDs  []int
	IsReady  bool
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
	Title       string
	Description string
	CardType    string
	Ranking     string
	URL         string
}

// GetNewCard :
func GetNewCard(state *State) Card {
	index := GetNewID(GetCardIDs(state))
	return Card{IDable: IDable{ID: index}, Info: CardInfo{}}
}

// GetNewPlayer :
func GetNewPlayer(state *State) Player {
	index := GetNewID(GetPlayerIDs(state))
	return Player{IDable: IDable{ID: index}, Info: PlayerInfo{Name: "New Player"}}
}

// GetNewGame :
func GetNewGame(state *State) Game {
	index := GetNewID(GetGameIDs(state))
	return Game{IDable: IDable{ID: index}, Info: GameInfo{Title: "New Game", Players: []PlayerGameLink{}, Battles: []Battle{}}}
}

// GetCardIDs :
func GetCardIDs(state *State) []IDable {
	var list []IDable
	for i := 0; i < len(state.Cards); i++ {
		list = append(list, state.Cards[i].IDable)
	}
	return list
}

// GetPlayerIDs :
func GetPlayerIDs(state *State) []IDable {
	var list []IDable
	for i := 0; i < len(state.Players); i++ {
		list = append(list, state.Players[i].IDable)
	}
	return list
}

// GetGameIDs :
func GetGameIDs(state *State) []IDable {
	var list []IDable
	for i := 0; i < len(state.Games); i++ {
		list = append(list, state.Games[i].IDable)
	}
	return list
}

// GetNewID :
func GetNewID(list []IDable) int {
	if len(list) == 0 {
		return 0
	}
	return list[len(list)-1].ID + 1
}

// GetIndexFromID :
func GetIndexFromID(list []IDable, id int) int {
	for i := 0; i < len(list); i++ {
		if list[i].ID == id {
			return i
		}
	}
	return -1
}

// GetNewBattle :
func GetNewBattle(player1 PlayerGameLink, player2 PlayerGameLink) Battle {
	battleSide1 := BattleSide{PlayerID: player1.PlayerID, Info: BattleSideInfo{Bench: []CardSlot{}, Hand: []CardSlot{}, TakePile: []CardSlot{}, ThrowPile: []CardSlot{}}}
	battleSide2 := BattleSide{PlayerID: player2.PlayerID, Info: BattleSideInfo{Bench: []CardSlot{}, Hand: []CardSlot{}, TakePile: []CardSlot{}, ThrowPile: []CardSlot{}}}
	for i := 0; i < 5; i++ {
		battleSide1.Info.Bench = append(battleSide1.Info.Bench, CardSlot{CardID: -1})
		battleSide2.Info.Bench = append(battleSide2.Info.Bench, CardSlot{CardID: -1})
	}

	for i := 0; i < len(player1.CardIDs); i++ {
		battleSide1.Info.TakePile = append(battleSide1.Info.TakePile, CardSlot{CardID: player1.CardIDs[i]})
	}

	for i := 0; i < len(player2.CardIDs); i++ {
		battleSide2.Info.TakePile = append(battleSide2.Info.TakePile, CardSlot{CardID: player2.CardIDs[i]})
	}

	return Battle{Sides: []BattleSide{battleSide1, battleSide2}, Turn: 0}
}

// DrawCards :
func DrawCards(state *State) []*Card {
	var drawnLegends int
	var drawnRares int
	var drawnMinions int
	var result []*Card

	legendsMax := 5
	raresMax := 10
	minionsMax := 20

	rand.Seed(time.Now().UnixNano())

	for drawnLegends < legendsMax || drawnRares < raresMax || drawnMinions < minionsMax {
		drawnCard := &state.Cards[rand.Intn(len(state.Cards))]
		if drawnLegends < legendsMax && drawnCard.Info.Ranking == "Legend" {
			drawnLegends++
			result = append(result, drawnCard)
		} else if drawnRares < raresMax && drawnCard.Info.Ranking == "Rare" {
			drawnRares++
			result = append(result, drawnCard)
		} else if drawnMinions < minionsMax && drawnCard.Info.Ranking == "Minion" {
			drawnMinions++
			result = append(result, drawnCard)
		}
	}

	return result
}

// StartGame :
func StartGame(state *State, game *Game) {
	for i := 0; i < len(game.Info.Players); i++ {
		player := &game.Info.Players[i]
		cards := DrawCards(state)

		for j := 0; j < len(cards); j++ {
			player.CardIDs = append(player.CardIDs, cards[j].ID)
		}
	}
	side0 := BattleSide{PlayerID: 0, Info: BattleSideInfo{Bench: []CardSlot{}, Hand: []CardSlot{}, ThrowPile: []CardSlot{}}}
	side1 := BattleSide{PlayerID: 2, Info: BattleSideInfo{Bench: []CardSlot{}, Hand: []CardSlot{}, ThrowPile: []CardSlot{}}}
	battle := Battle{Sides: []BattleSide{side0, side1}}
	for i, _ := range battle.Sides {
		side := &battle.Sides[i]
		for _, id := range game.Info.Players[side.PlayerID].CardIDs {
			side.Info.TakePile = append(side.Info.TakePile, CardSlot{CardID: id})
		}
	}
	game.Info.Battles = append(game.Info.Battles, battle)
}

func readJSON(state *State) {
	jsonFile, err := os.Open("./cards/JSONCards.txt")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var infos []CardInfo
	json.Unmarshal(byteValue, &infos)

	for i := 0; i < len(infos); i++ {
		card := GetNewCard(state)
		card.Info = infos[i]
		state.Cards = append(state.Cards, card)
	}
}

// StartBattle :
func StartBattle(battle *Battle) {
	for i, _ := range battle.Sides {
		side := &battle.Sides[i]
		DrawCard(&side.Info)
		DrawCard(&side.Info)
		DrawCard(&side.Info)
		DrawCard(&side.Info)
		DrawCard(&side.Info)
	}
}

// BattleNextPlayer :
func BattleStartTurn(battle *Battle) {

}

// BattleEndTurn :
func BattleEndTurn(battle *Battle) {
	side := &battle.Sides[battle.Turn]
	DrawCard(&side.Info)
	battle.Turn = (battle.Turn + 1) % len(battle.Sides)
}

// DrawCard :
func DrawCard(side *BattleSideInfo) {
	takePile := &side.TakePile
	if len(*takePile) > 0 {
		drawCard := (*takePile)[0]
		side.TakePile = (*takePile)[1:]
		side.Hand = append(side.Hand, drawCard)
	}
}

// GetExampleState :
func GetExampleState() State {
	var state = State{Games: []Game{}, Players: []Player{}, Cards: []Card{}}

	readJSON(&state)

	player1 := GetNewPlayer(&state)
	player1.Info.Name = "dj leo"

	state.Players = append(state.Players, player1)

	player2 := GetNewPlayer(&state)
	player2.Info.Name = "funky kong"

	state.Players = append(state.Players, player2)

	player3 := GetNewPlayer(&state)
	player3.Info.Name = "richard"

	state.Players = append(state.Players, player3)

	game1 := GetNewGame(&state)
	game1.Info.Title = "Cool game!"

	var deck []int

	link1 := PlayerGameLink{PlayerID: 0, CardIDs: deck, IsReady: true}
	link2 := PlayerGameLink{PlayerID: 1, CardIDs: deck, IsReady: true}
	link3 := PlayerGameLink{PlayerID: 2, CardIDs: deck}

	game1.Info.Players = append(game1.Info.Players, link1, link2, link3)
	state.Games = append(state.Games, game1)

	return state
}
