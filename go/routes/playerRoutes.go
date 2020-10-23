package routes

import (
	"net/http"
	"strconv"

	"../lio"
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

	if id != -1 {
		lio.HandleGETResponse(w, state.Players[id])
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
		index = id
	}

	if index != -1 {
		lio.HandleGETResponse(w, state.Players[index])
	} else {
		lio.HandleGETResponse(w, "")
	}
}
