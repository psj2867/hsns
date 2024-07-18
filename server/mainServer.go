package server

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/router"
	"github.com/psj2867/hsns/server/middleware"
)

func InitServer() *gin.Engine {
	r := gin.Default()
	r.Use(corsMiddleware())
	r.Use(middleware.GlobalAuthMiddleware())
	router.SetRouter(r)
	return r
}
func corsMiddleware() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AddAllowHeaders("*")
	config.AllowAllOrigins = true
	return cors.New(config)
}
func DeferServer(r *gin.Engine) {
	config.DeferDb()
}

func Run(addr ...string) {
	s := InitServer()
	defer DeferServer(s)
	s.Run(addr...)
}
