package router

import "github.com/gin-gonic/gin"

func SetRouter(r *gin.Engine) {
	r.GET("", func(ctx *gin.Context) { ctx.Writer.WriteString("/") })
	userRouter{}.setRouter(r.Group("/user"))
	contents{}.setRouter(r.Group("/contents"))
}
