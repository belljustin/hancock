package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belljustin/hancock/models"
	_ "github.com/belljustin/hancock/models/mem"
	_ "github.com/belljustin/hancock/models/postgres"
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

type appHandler func(http.ResponseWriter, *http.Request) error

func (fn appHandler) Handle(w http.ResponseWriter, req *http.Request) {
	code := 500
	if err := fn(w, req); err != nil {
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

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func NewRouter(c *Config) http.Handler {
	router := gin.Default()
	router.GET("/ping", ping)

	keys, err := models.Open(c.StorageType, c.StorageConfig)
	if err != nil {
		panic(err)
	}
	algs := map[string]Alg{
		"rsa": &Rsa{},
	}

	RegisterKeyHandlers(router, keys, algs)

	return router
}
