package routes

import (
	"net/http"
	"strconv"

	"../lio"
	"../model"
	"github.com/julienschmidt/httprouter"
)

// SetCardRoutes :
func SetCardRoutes() {
	router.GET("/card/:id", GETCard)
	router.GET("/cards/:game-id", GETCards)
}

// GETCard :
func GETCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := lio.GetIntParam(ps, "id")

	if index != -1 {
		lio.HandleGETResponse(w, state.Cards[index])
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETCards :
func GETCards(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	if gameID >= 0 && gameID < len(state.Games) {
		game := state.Games[gameID]

		playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
		cardIDs := game.Info.Players[playerID].CardIDs

		var cards []model.Card
		for i := 0; i < len(cardIDs); i++ {
			cards = append(cards, state.Cards[cardIDs[i]])
		}

		lio.HandleGETResponse(w, cards)
	}
}
