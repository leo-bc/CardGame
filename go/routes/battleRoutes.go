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
	router.POST("/battle-start/:game-id/:battle-id/", POSTBattleStart)

	router.POST("/battle-end-turn/:game-id/:battle-id", POSTBattleEndTurn)
	router.POST("/battle-play-card/:game-id/:battle-id/:hand-id", POSTBattlePlayCard)
	router.POST("/battle-attack/:game-id/:battle-id/", POSTBattleAttack)

	router.GET("/battle-updated/:game-id/:battle-id/:update-id", GETBattleUpdated)
}

// GETBattleUpdated :
func GETBattleUpdated(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")
	updateID := lio.GetIntParam(ps, "update-id")

	if gameID != -1 && battleID != -1 && len(state.Games) > gameID && len(state.Games[gameID].Info.Battles) > battleID {
		updateInfo := state.Games[gameID].Info.Battles[battleID].TriggerUpdate(updateID)
		lio.HandleGETResponse(w, updateInfo)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// POSTBattleEndTurn :
func POSTBattleEndTurn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		battle := &state.Games[gameID].Info.Battles[battleID]
		model.BattleEndTurn(battle)
		battle.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// POSTBattlePlayCard :
func POSTBattlePlayCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")
	handID := lio.GetIntParam(ps, "hand-id")
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))

	if gameID != -1 && battleID != -1 && handID != -1 {
		var side *model.BattleSide
		battle := &state.Games[gameID].Info.Battles[battleID]
		if battle.Sides[0].PlayerID == playerID {
			side = &battle.Sides[0]
		} else {
			side = &battle.Sides[1]
		}
		model.PlayCard(side, handID)
		battle.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// POSTBattleAttack :
func POSTBattleAttack(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
	var info model.AttackInfo
	lio.DecodePOSTBody(r, &info)

	if gameID != -1 && battleID != -1 {
		var playerSide, opponentSide *model.BattleSide
		battle := &state.Games[gameID].Info.Battles[battleID]
		if battle.Sides[0].PlayerID == playerID {
			playerSide = &battle.Sides[0]
			opponentSide = &battle.Sides[1]
		} else {
			playerSide = &battle.Sides[1]
			opponentSide = &battle.Sides[0]
		}

		model.AttackCard(state, battle, playerSide, opponentSide, info, vm)
		battle.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// POSTBattleStart :
func POSTBattleStart(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		battle := &state.Games[gameID].Info.Battles[battleID]
		battle.Info.IsStarted = true
		model.StartBattle(battle)
		battle.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

func createClientBattleSides(sides *[]model.BattleSide, turn int, playerID int) (*[]model.ClientBattleSide, []int) {
	allCards := []int{}
	allSides := []model.ClientBattleSide{}
	for i := range *sides {
		side := (*sides)[i]
		slots := make(map[string][]model.CardSlot)
		for title, container := range side.Info.Cards {
			slots[title] = []model.CardSlot{}
			if (title == "Hand" || title == "TakePile") && side.PlayerID != playerID {
				for i := 0; i < len(container.Slots); i++ {
					slots[title] = append(slots[title], model.CardSlot{CardID: -1})
				}
			} else {
				for _, slot := range container.Slots {
					slots[title] = append(slots[title], slot)
					allCards = append(allCards, slot.CardID)
				}
			}
		}

		clientSide := model.ClientBattleSide{IsPlayer: side.PlayerID == playerID, IsTurn: turn == i, Cards: slots}
		allSides = append(allSides, clientSide)
	}
	return &allSides, allCards
}

// GETBattle :
func GETBattle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := lio.GetIntParam(ps, "game-id")
	battleID := lio.GetIntParam(ps, "battle-id")
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))

	if gameID != -1 && battleID != -1 && len(state.Games) > gameID && len(state.Games[gameID].Info.Battles) > battleID {
		battle := state.Games[gameID].Info.Battles[battleID]

		allSides, allCards := createClientBattleSides(&battle.Sides, battle.Info.Turn, playerID)

		cards := make(map[int]model.Card)
		for _, cardID := range allCards {
			cards[cardID] = state.Cards[cardID]
		}

		clientBattle := model.ClientBattle{Sides: *allSides, Cards: cards, Info: battle.Info}
		lio.HandleGETResponse(w, clientBattle)
	} else {
		lio.HandleGETResponse(w, "")
	}
}
