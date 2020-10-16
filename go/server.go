package main

import (
	"log"
	"net/http"

	"./model"
	"./routes"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
)

var state *model.State
var cookieEncoder *securecookie.SecureCookie

func main() {
	s := model.GetExampleState()
	state = &s
	createRouter()
}

func createRouter() {
	var hashKey = []byte("15989999955333994")
	var blockKey = []byte("1234567812345678")
	cookieEncoder = securecookie.New(hashKey, blockKey)

	router := httprouter.New()
	routes.SetRoutes(router, cookieEncoder, state)

	// for debugging purposes
	router.ServeFiles("/website/*filepath", http.Dir("./website/"))

	log.Printf("starting server at http://localhost:10001\n")
	log.Fatal(http.ListenAndServe(":10001", router))
}
