package router

import "github.com/gin-gonic/gin"

func SetRouter(r *gin.Engine) {
	user{}.setRouter(r.Group("/user"))
	contents{}.setRouter(r.Group("/contents"))
}
