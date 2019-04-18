package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belljustin/hancock/models"
	_ "github.com/belljustin/hancock/models/mem"
	_ "github.com/belljustin/hancock/models/postgres"
)

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

	registerKeyHandlers(router, keys, algs)

	return router
}
