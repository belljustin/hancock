package server

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/belljustin/hancock/key"
)

func ping(c *gin.Context) {
	c.String(http.StatusOK, "Pong")
}

func Run(port int, s key.Storage) error {
	router := gin.Default()

	router.GET("/ping", ping)
	registerKeyHandlers(router, s)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), router)
}
