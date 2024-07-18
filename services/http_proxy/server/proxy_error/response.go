package proxy_error

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HTTPErrorBody struct {
	Code, Message string
}

func AbortWithErrorResponse(c *gin.Context, e ProxyError, message string) {
	c.AbortWithStatusJSON(e.HttpStatusCode(), HTTPErrorBody{
		Code:    e.String(),
		Message: message,
	})
}

func AbortWithErrorResponsef(c *gin.Context, e ProxyError, format string, args ...interface{}) {
	c.AbortWithStatusJSON(e.HttpStatusCode(), HTTPErrorBody{
		Code:    e.String(),
		Message: fmt.Sprintf(format, args...),
	})
}
