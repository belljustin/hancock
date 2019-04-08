package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/fakes"
)

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Pong")
}

func main() {

	router := httprouter.New()
	router.GET("/ping", ping)

	keys := fakes.Keys{}
	RegisterKeyHandlers(router, keys)

	http.ListenAndServe(":8000", router)
}
