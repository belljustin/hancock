package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belljustin/hancock/key"
	_ "github.com/belljustin/hancock/key/mem"
	_ "github.com/belljustin/hancock/key/postgres"
)

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func NewRouter(c *Config) http.Handler {
	router := gin.Default()
	router.GET("/ping", ping)

	keyStorage, err := key.Open(c.StorageType, c.StorageConfig)
	if err != nil {
		panic(err)
	}

	registerKeyHandlers(router, keyStorage)

	return router
}
