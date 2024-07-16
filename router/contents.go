package router

import (
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/server/middleware"
)

type contents struct{}

func (t contents) setRouter(group *gin.RouterGroup) {
	group.GET("", t.getContent)
	uploadGroup := group.Group("/upload", middleware.ShouldAuthUnauth())
	{
		uploadGroup.POST("", t.upload)
		uploadGroup.POST("/success", t.uploadSuccess)
		uploadGroup.GET("/fail", t.uploadFail)
	}
}
func (t *contents) getContent(c *gin.Context) {
}
