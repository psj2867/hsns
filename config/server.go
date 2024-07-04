package config

import (
	"github.com/gin-gonic/gin"
)

func InitServer() *gin.Engine {
	r := gin.Default()
	return r
}
