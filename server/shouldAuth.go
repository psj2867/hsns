package server

import (
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/config"
)

func ShouldAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Get(config.JwtUser)
		c.G
	}
}
