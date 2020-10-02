package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"./lio"
	"./model"
	"github.com/julienschmidt/httprouter"
)

var state model.State

// GETIndex : redirect to the actual homepage for ease of use
func GETIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusTemporaryRedirect) // 307
	http.Redirect(w, r, "/static/", 307)
}

// GETState :
func GETState(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	json.NewEncoder(w).Encode(state)
}

// GETCard :
func GETCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := getIndexFromID(getCardIDs(), getIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Cards[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETGame :
func GETGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := getIndexFromID(getGameIDs(), getIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Games[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETPlayer :
func GETPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := getIndexFromID(getPlayerIDs(), getIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Players[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETPlayers :
func GETPlayers(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lio.HandleGETResponse(w, state.Players)
}

// GETCards :
func GETCards(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lio.HandleGETResponse(w, state.Cards)
}

// POSTCard :
func POSTCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var info model.CardInfo
	lio.DecodePOSTBody(r, &info)
	cardID := getIntParam(ps, "id")
	index := getIndexFromID(getCardIDs(), cardID)

	if index != -1 {
		state.Cards[index].Info = info
	}
	lio.HandlePOSTResponse(w)
}

// GETCreateCard :
func GETCreateCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	newCard := getNewCard()
	state.Cards = append(state.Cards, newCard)
	lio.HandleGETResponse(w, newCard.ID)
}

// POSTRemoveCard :
func POSTRemoveCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cardID := getIntParam(ps, "id")
	var cards []model.Card
	for i := 0; i < len(state.Cards); i++ {
		if state.Cards[i].ID != cardID {
			cards = append(cards, state.Cards[i])
		}
	}
	state.Cards = cards
	lio.HandlePOSTResponse(w)
}

func getIntParam(ps httprouter.Params, name string) int {
	value, _ := strconv.Atoi(ps.ByName(name))
	return value
}

func getCardIDs() []model.IDable {
	var list []model.IDable
	for i := 0; i < len(state.Cards); i++ {
		list = append(list, state.Cards[i].IDable)
	}
	return list
}

func getPlayerIDs() []model.IDable {
	var list []model.IDable
	for i := 0; i < len(state.Players); i++ {
		list = append(list, state.Players[i].IDable)
	}
	return list
}

func getGameIDs() []model.IDable {
	var list []model.IDable
	for i := 0; i < len(state.Games); i++ {
		list = append(list, state.Games[i].IDable)
	}
	return list
}

func getNewCard() model.Card {
	index := getNewID(getCardIDs())
	return model.Card{IDable: model.IDable{ID: index}, Info: model.CardInfo{Title: "New Card", HP: 10}}
}

func getNewPlayer() model.Player {
	index := getNewID(getPlayerIDs())
	return model.Player{IDable: model.IDable{ID: index}, Info: model.PlayerInfo{Name: "New Player"}}
}

func getNewGame() model.Game {
	index := getNewID(getGameIDs())
	return model.Game{IDable: model.IDable{ID: index}, Info: model.GameInfo{Title: "New Game", Players: []model.PlayerGameLink{}}}
}

func getNewID(list []model.IDable) int {
	if len(list) == 0 {
		return 0
	}
	return list[len(list)-1].ID + 1
}

func getIndexFromID(list []model.IDable, id int) int {
	for i := 0; i < len(list); i++ {
		if list[i].ID == id {
			return i
		}
	}
	return -1
}

func main() {
	state = model.State{Games: []model.Game{}, Players: []model.Player{}, Cards: []model.Card{}}

	card1 := getNewCard()
	card1.Info.Title = "Leo"
	card1.Info.HP = 20

	state.Cards = append(state.Cards, card1)

	card2 := getNewCard()
	card2.Info.Title = "Jop"

	state.Cards = append(state.Cards, card2)

	card3 := getNewCard()
	card3.Info.Title = "Ark"

	state.Cards = append(state.Cards, card3)

	player1 := getNewPlayer()
	player1.Info.Name = "dj leo"

	state.Players = append(state.Players, player1)

	player2 := getNewPlayer()
	player2.Info.Name = "funky kong"

	state.Players = append(state.Players, player2)

	player3 := getNewPlayer()
	player3.Info.Name = "richard"

	state.Players = append(state.Players, player3)

	game1 := getNewGame()
	game1.Info.Title = "Cool game!"

	link1 := model.PlayerGameLink{PlayerID: 0, CardIDs: []int{0}}
	link2 := model.PlayerGameLink{PlayerID: 1, CardIDs: []int{1, 2}}
	link3 := model.PlayerGameLink{PlayerID: 2, CardIDs: []int{}}

	game1.Info.Players = append(game1.Info.Players, link1, link2, link3)
	state.Games = append(state.Games, game1)

	router := httprouter.New()
	router.GET("/", GETIndex)
	router.GET("/state", GETState)
	router.GET("/player/:id", GETPlayer)
	router.GET("/players/", GETPlayers)
	router.GET("/game/:id", GETGame)
	router.GET("/card/:id", GETCard)
	router.GET("/cards/", GETCards)
	router.POST("/card/:id", POSTCard)
	router.GET("/create-card/", GETCreateCard)
	router.POST("/remove-card/:id", POSTRemoveCard)

	// for debugging purposes
	router.ServeFiles("/static/*filepath", http.Dir("./static/"))

	log.Printf("starting server at http://localhost:10001\n")
	log.Fatal(http.ListenAndServe(":10001", router))
}
