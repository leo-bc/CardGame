package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"./model"
	"./routes"
	"github.com/gorilla/securecookie"
	"github.com/julienschmidt/httprouter"
	"github.com/robertkrimen/otto"
)

var state *model.State
var cookieEncoder *securecookie.SecureCookie

func main() {
	s := model.GetExampleState(false)
	state = &s
	createRouter()
}

func createRouter() {
	var hashKey = []byte("15989999955333994")
	var blockKey = []byte("1234567812345678")
	cookieEncoder = securecookie.New(hashKey, blockKey)

	vm := otto.New()

	vm.Set("flipCoin", func(call otto.FunctionCall) otto.Value {
		rand.Seed(time.Now().UnixNano())
		throw := rand.Intn(2)
		value, _ := otto.ToValue(throw == 1)
		return value
	})

	vm.Run(`
		function toJSON(any) {
			return JSON.stringify(any, null, 4)
		}

		var CONFUSED = "Confused";
		var POISONED = "Poisoned";
		var PARALYZED = "Paralyzed";
		var ASLEEP = "Asleep";
	`)

	router := httprouter.New()
	routes.SetRoutes(router, cookieEncoder, state, vm)
	router.ServeFiles("/website/*filepath", http.Dir("./website/"))

	log.Printf("starting server at http://localhost:10001\n")
	log.Fatal(http.ListenAndServe(":10001", router))
}
