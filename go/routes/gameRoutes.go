package routes

import (
	"net/http"
	"strconv"

	"../lio"
	"../model"
	"github.com/julienschmidt/httprouter"
)

// SetGameRoutes :
func SetGameRoutes() {
	router.GET("/game/:id", GETGame)

	router.GET("/game-updated/:id", GETGameUpdated)

	router.POST("/join-game/:id", POSTJoinGame)
	router.POST("/leave-game/:id", POSTLeaveGame)

	router.POST("/game-set-ready/:id", POSTGameSetReady)
}

// GETGame :
func GETGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), lio.GetIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Games[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETGameUpdated :
func GETGameUpdated(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), lio.GetIntParam(ps, "id"))

	if index != -1 {
		isUpdated := state.Games[index].TriggerUpdate()
		lio.HandleGETResponse(w, isUpdated)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// POSTJoinGame :
func POSTJoinGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), lio.GetIntParam(ps, "id"))
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
	if index != -1 && playerID != -1 {
		link := model.PlayerGameLink{PlayerID: playerID, CardIDs: []int{}}
		state.Games[index].Info.Players = append(state.Games[index].Info.Players, link)
		state.Games[index].SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// POSTLeaveGame :
func POSTLeaveGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), lio.GetIntParam(ps, "id"))
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
	if index != -1 && playerID != -1 {
		var links []model.PlayerGameLink
		players := state.Games[index].Info.Players
		for i := 0; i < len(players); i++ {
			if players[i].PlayerID != playerID {
				links = append(links, players[i])
			}
		}
		state.Games[index].Info.Players = links
		state.Games[index].SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// POSTGameSetReady :
func POSTGameSetReady(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := model.GetIndexFromID(model.GetGameIDs(state), lio.GetIntParam(ps, "id"))
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
	if gameID != -1 && playerID != -1 {
		game := &state.Games[gameID]
		player := &game.Info.Players[playerID]
		player.IsReady = !player.IsReady

		if player.IsReady {
			everyoneReady := true
			for i := 0; i < len(game.Info.Players); i++ {
				if !game.Info.Players[i].IsReady {
					everyoneReady = false
					break
				}
			}
			if everyoneReady {
				game.Info.IsStarted = true
				model.StartGame(state, game)
			}
		}

		game.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}
