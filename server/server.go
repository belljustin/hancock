package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/belljustin/hancock/models"
	_ "github.com/belljustin/hancock/models/mem"
)

type httpError struct {
	Code    int
	Message string
}

func (err *httpError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func newInternalServerError(err error) *httpError {
	return &httpError{
		http.StatusInternalServerError,
		err.Error(),
	}
}

type appHandler func(http.ResponseWriter, *http.Request, httprouter.Params) error

func (fn appHandler) Handle(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	code := 500
	if err := fn(w, req, ps); err != nil {
		if err, ok := err.(*httpError); ok {
			code = err.Code
		}

		if code > 500 {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", code)
		} else {
			http.Error(w, err.Error(), code)
		}
	}
}

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	io.WriteString(w, "Pong")
}

func NewRouter() http.Handler {
	router := httprouter.New()
	router.GET("/ping", ping)

	keys, err := models.Open("mem")
	if err != nil {
		panic(err)
	}
	algs := map[string]Alg{
		"rsa": &Rsa{},
	}

	RegisterKeyHandlers(router, keys, algs)

	return router
}
