package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrHandler func(c *gin.Context)

func ShouldAuthMiddleware(errHanldler ErrHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		ok := HasAuth(c)
		if ok {
			return
		}
		errHanldler(c)
	}
}

func ShouldAuthUnauth() gin.HandlerFunc {
	return ShouldAuthMiddleware(AbortError(http.StatusUnauthorized))
}
func AbortError(code int) ErrHandler {
	return func(c *gin.Context) {
		c.AbortWithStatus(code)
	}
}
