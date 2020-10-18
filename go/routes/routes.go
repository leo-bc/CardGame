package routes

import (
	"../model"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
)

var router *httprouter.Router
var cookieEncoder *securecookie.SecureCookie
var state *model.State

// SetRoutes :
func SetRoutes(r *httprouter.Router, c *securecookie.SecureCookie, s *model.State) {
	router = r
	cookieEncoder = c
	state = s

	SetPlayerRoutes()
	SetCardRoutes()
	SetGameRoutes()
	SetBattleRoutes()
}
