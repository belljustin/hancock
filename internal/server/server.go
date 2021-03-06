package server

import (
	_ "encoding/json" // for tagging json structs
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belljustin/hancock/key"
)

// Config provides settings for the hancock REST server.
type Config struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

// Run a hancock REST server using s as the backend `key.Storage`.
func Run(port int, s key.Storage) error {
	router := gin.Default()

	router.GET("/ping", ping)
	registerKeyHandlers(router, s)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
