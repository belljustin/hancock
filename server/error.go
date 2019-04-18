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
	if err, ok := err.(*httpError); ok {
		c.Error(err)
		c.AbortWithStatusJSON(err.Code, &err)
	} else {
		panic(err)
	}
}
