package main

import (
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/models"
	_ "github.com/belljustin/hancock/models/mem"
)

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Pong")
}

func main() {

	router := httprouter.New()
	router.GET("/ping", ping)

	keys, ok := models.Open("mem")
	if !ok {
		panic("Could not initialize key storage")
	}
	algs := map[string]Alg{
		"rsa": &Rsa{},
	}

	RegisterKeyHandlers(router, keys, algs)

	err := http.ListenAndServe(":8000", router)
	panic(err)
}
