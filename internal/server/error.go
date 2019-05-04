package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type httpError struct {
	Code    int
	Message string
}

func (err *httpError) Error() string {
	return fmt.Sprintf("%d: %s", err.Code, err.Message)
}

func handleError(c *gin.Context, err error) {
	if herr, ok := err.(*httpError); ok {
		c.Error(herr)
		c.AbortWithStatusJSON(herr.Code, &herr)
	} else {
		panic(err)
	}
}
