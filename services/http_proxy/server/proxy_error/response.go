package proxy_error

import (
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
