package router

import "github.com/gin-gonic/gin"

func SetRouter(r *gin.Engine) {
	test := r.Group("/test")
	{
		test.GET("/asdf", func(ctx *gin.Context) {
			ctx.Writer.WriteString("asdfasdf")
		})
	}
}
