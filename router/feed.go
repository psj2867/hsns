package router

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/models"
	"github.com/psj2867/hsns/server/middleware"
)

type feed struct{}

func (t feed) setRouter(group *gin.RouterGroup) {
	group.GET("/", t.getFeed)
	group.GET("/main", t.getMain, middleware.ShouldAuthUnauth())
}
func (t *feed) getFeed(c *gin.Context) {

}
func (t *feed) getMain(c *gin.Context) {
	user := models.GetUserInfoInContext(c)
	contents := models.Contents{}
	err := contents.GetByUser(user.Id)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return
	}
	fmt.Printf("contents: %v\n", contents)

}
