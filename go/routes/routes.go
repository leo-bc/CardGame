package routes

import (
	"net/http"

	"../lio"
	"../model"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"github.com/robertkrimen/otto"
)

var router *httprouter.Router
var cookieEncoder *securecookie.SecureCookie
var state *model.State
var vm *otto.Otto

// SetRoutes :
func SetRoutes(r *httprouter.Router, c *securecookie.SecureCookie, s *model.State, v *otto.Otto) {
	router = r
	cookieEncoder = c
	state = s
	vm = v
	router.GET("/to-json/:with-cards", GETJSON)

	SetPlayerRoutes()
	SetCardRoutes()
	SetGameRoutes()
	SetBattleRoutes()
}

// StateJSON :
type StateJSON struct {
	Players []model.Player
	Games   []model.Game
}

// GETJSON :
func GETJSON(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	if lio.GetIntParam(ps, "with-cards") == 0 {
		lio.HandleGETResponse(w, StateJSON{state.Players, state.Games})
	} else {
		lio.HandleGETResponse(w, state)
	}
}
