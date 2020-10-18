package routes

import (
	"net/http"
	"strconv"

	"../lio"
	"../model"
	"github.com/julienschmidt/httprouter"
)

// SetBattleRoutes :
func SetBattleRoutes() {
	router.GET("/battle/:game-id/:battle-id", GETBattle)
	router.GET("/battle-side/:game-id/:battle-id/:side-id", GETBattleSide)
	router.POST("/battle-start/:game-id/:battle-id/", POSTBattleStart)

	router.POST("/battle-end-turn/:game-id/:battle-id", POSTBattleEndTurn)

}

// POSTBattleEndTurn :
func POSTBattleEndTurn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		battle := &state.Games[gameID].Info.Battles[battleID]
		model.BattleEndTurn(battle)
	}
	lio.HandlePOSTResponse(w)
}

// POSTBattleStart :
func POSTBattleStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		battle := &state.Games[gameID].Info.Battles[battleID]
		battle.IsStarted = true
		model.StartBattle(battle)
	}
	lio.HandlePOSTResponse(w)
}

// GETBattle :
func GETBattle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		lio.HandleGETResponse(w, state.Games[gameID].Info.Battles[battleID])
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETBattleSide :
func GETBattleSide(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")
	sideID := lio.GetIntParam(ps, "side-id")
	if gameID != -1 && battleID != -1 {
		battle := state.Games[gameID].Info.Battles[battleID]
		side := battle.Sides[sideID]
		playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
		var sideInfo model.PlayerBattleSideInfo
		if side.PlayerID == playerID {
			sideInfo = model.PlayerBattleSideInfo{IsPlayer: true, IsTurn: battle.Turn == sideID, Info: side.Info}
		} else {
			sideInfo = model.PlayerBattleSideInfo{IsPlayer: false, IsTurn: battle.Turn == sideID, Info: side.Info}
		}
		lio.HandleGETResponse(w, sideInfo)
	} else {
		lio.HandleGETResponse(w, "")
	}
}
