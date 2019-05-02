package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/belljustin/hancock/key/postgres"
)

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func NewRouter(c *Config) http.Handler {
	router := gin.Default()
	router.GET("/ping", ping)

	keyStorage := db.KeyStorage{}
	err := keyStorage.Open(c.StorageConfig)
	if err != nil {
		panic(err)
	}

	registerKeyHandlers(router, &keyStorage)

	return router
}
