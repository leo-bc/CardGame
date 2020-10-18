package routes

import (
	"net/http"
	"strconv"

	"../lio"
	"../model"
	"github.com/julienschmidt/httprouter"
)

// SetPlayerRoutes :
func SetPlayerRoutes() {
	router.GET("/player/:id", GETPlayer)
	router.GET("/players/", GETPlayers)

	router.POST("/select-player/:id", POSTSelectPlayer)
	router.GET("/current-player/", GETCurrentPlayer)
}

// GETPlayer :
func GETPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := lio.GetIntParam(ps, "id")
	index := model.GetIndexFromID(model.GetPlayerIDs(state), id)

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

// POSTSelectPlayer :
func POSTSelectPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	playerID := lio.GetIntParam(ps, "id")
	idString := strconv.Itoa(playerID)
	lio.SetCookie(cookieEncoder, w, "player-id", idString)
	lio.HandlePOSTResponse(w)
}

// GETCurrentPlayer :
func GETCurrentPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	cookieResult := lio.ReadCookie(cookieEncoder, r, "player-id")
	index := -1
	if cookieResult != "" {
		id, _ := strconv.Atoi(cookieResult)
		index = model.GetIndexFromID(model.GetPlayerIDs(state), id)
	}

	if index != -1 {
		lio.HandleGETResponse(w, state.Players[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}
