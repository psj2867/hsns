package server

import (
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/router"
)

func InitServer() *gin.Engine {
	r := gin.Default()
	r.Use(globalAuthMiddleware())
	router.SetRouter(r)
	return r
}
