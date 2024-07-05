package router

import (
	"github.com/gin-gonic/gin"
)

func InitServer() *gin.Engine {
	r := gin.Default()
	SetRouter(r)
	return r
}
