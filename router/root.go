package router

import "github.com/gin-gonic/gin"

func SetRouter(r *gin.Engine) {
	user{}.setUserRouter(r.Group("/user"))

}
