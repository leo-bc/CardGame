package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"./lio"
	"./model"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
)

var state model.State
var cookieEncoder *securecookie.SecureCookie

func main() {
	state = model.GetExampleState()
	createRouter()
}

func createRouter() {
	var hashKey = []byte("15989999955333994")
	var blockKey = []byte("1234567812345678")
	cookieEncoder = securecookie.New(hashKey, blockKey)

	router := httprouter.New()
	router.GET("/", GETIndex)
	router.GET("/state", GETState)
	router.GET("/player/:id", GETPlayer)
	router.GET("/current-player/", GETCurrentPlayer)
	router.GET("/players/", GETPlayers)
	router.POST("/select-player/:id", POSTSelectPlayer)
	router.GET("/selected-player/", GETSelectedPlayer)
	router.GET("/game/:id", GETGame)
	router.GET("/game-updated/:id", GETGameUpdated)
	router.POST("/join-game/:id", POSTJoinGame)
	router.POST("/leave-game/:id", POSTLeaveGame)
	router.POST("/set-ready/:id", POSTSetReady)
	router.GET("/battle/:game-id/:battle-id", GETBattle)
	router.GET("/battle-side/:game-id/:battle-id/:side-id", GETBattleSide)
	router.GET("/card/:id", GETCard)
	router.GET("/cards/:game-id", GETCards)
	router.POST("/draw-card/:game-id/:battle-id", POSTDrawCard)

	router.POST("/card/:id", POSTCard)
	router.GET("/create-card/", GETCreateCard)
	router.POST("/remove-card/:id", POSTRemoveCard)

	// for debugging purposes
	router.ServeFiles("/website/*filepath", http.Dir("../website/"))

	log.Printf("starting server at http://localhost:10001\n")
	log.Fatal(http.ListenAndServe(":10001", router))
}

// POSTSelectPlayer :
func POSTSelectPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	playerID := getIntParam(ps, "id")
	idString := strconv.Itoa(playerID)
	lio.SetCookie(cookieEncoder, w, "player-id", idString)
	lio.HandlePOSTResponse(w)
}

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
	index := model.GetIndexFromID(model.GetCardIDs(state), getIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Cards[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETGame :
func GETGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "id"))

	if index != -1 {
		lio.HandleGETResponse(w, state.Games[index].Info)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETGameUpdated :
func GETGameUpdated(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "id"))

	if index != -1 {
		isUpdated := state.Games[index].TriggerUpdate()
		lio.HandleGETResponse(w, isUpdated)
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// POSTJoinGame :
func POSTJoinGame(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	index := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "id"))
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
	index := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "id"))
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

// POSTSetReady :
func POSTSetReady(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "id"))
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
				model.StartGame(&state, game)
			}
		}

		game.SetUpdated()
	}
	lio.HandlePOSTResponse(w)
}

// GETBattle :
func GETBattle(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := getIntParam(ps, "game-id")
	battleID := getIntParam(ps, "battle-id")

	if gameID != -1 && battleID != -1 {
		lio.HandleGETResponse(w, state.Games[gameID].Info.Battles[battleID])
	} else {
		lio.HandleGETResponse(w, "")
	}
}

// GETBattleSide :
func GETBattleSide(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := getIntParam(ps, "game-id")
	battleID := getIntParam(ps, "battle-id")
	sideID := getIntParam(ps, "side-id")
	if gameID != -1 && battleID != -1 {
		battle := state.Games[gameID].Info.Battles[battleID]
		var side model.BattleSide
		if sideID == 0 {
			side = battle.FirstSide
		} else {
			side = battle.SecondSide
		}
		playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
		var sideInfo model.PlayerBattleSideInfo
		if side.PlayerID == playerID {
			sideInfo = model.PlayerBattleSideInfo{IsPlayer: true, Info: side.Info}
		} else {
			sideInfo = model.PlayerBattleSideInfo{IsPlayer: false, Info: side.Info}
		}
		lio.HandleGETResponse(w, sideInfo)
	} else {
		lio.HandleGETResponse(w, "")
	}
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

// GETPlayer :
func GETPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id := getIntParam(ps, "id")
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

// GETCards :
func GETCards(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	gameID := getIntParam(ps, "game-id")
	if gameID >= 0 && gameID < len(state.Games) {
		game := state.Games[gameID]

		playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
		playerIndex := model.GetIndexFromID(model.GetPlayerIDs(state), playerID)
		cardIDs := game.Info.Players[playerIndex].CardIDs

		var cards []model.Card
		for i := 0; i < len(cardIDs); i++ {
			cardID := model.GetIndexFromID(model.GetCardIDs(state), cardIDs[i])
			cards = append(cards, state.Cards[cardID])
		}

		lio.HandleGETResponse(w, cards)
	}
}

// POSTCard :
func POSTCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var info model.CardInfo
	lio.DecodePOSTBody(r, &info)
	cardID := getIntParam(ps, "id")
	index := model.GetIndexFromID(model.GetCardIDs(state), cardID)

	if index != -1 {
		state.Cards[index].Info = info
	}
	lio.HandlePOSTResponse(w)
}

// POSTDrawCard :
func POSTDrawCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	playerID, _ := strconv.Atoi(lio.ReadCookie(cookieEncoder, r, "player-id"))
	playerIndex := model.GetIndexFromID(model.GetPlayerIDs(state), playerID)
	gameIndex := model.GetIndexFromID(model.GetGameIDs(state), getIntParam(ps, "game-id"))
	battleID := getIntParam(ps, "battle-id")

	battle := state.Games[gameIndex].Info.Battles[battleID]
	if battle.FirstSide.PlayerID == playerIndex {
		model.DrawCard(&battle.FirstSide.Info)
	} else if battle.SecondSide.PlayerID == playerIndex {
		model.DrawCard(&battle.SecondSide.Info)
	}
	state.Games[gameIndex].Info.Battles[battleID] = battle
	lio.HandlePOSTResponse(w)
}

// GETCreateCard :
func GETCreateCard(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	newCard := model.GetNewCard(state)
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

// GETSelectedPlayer :
func GETSelectedPlayer(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	lio.HandleGETResponse(w, lio.ReadCookie(cookieEncoder, r, "player-id"))
}

func getIntParam(ps httprouter.Params, name string) int {
	value, _ := strconv.Atoi(ps.ByName(name))
	return value
}
