package router

import (
	"github.com/gin-gonic/gin"
	"github.com/guregu/null/v5"
	"github.com/psj2867/hsns/models"
)

type user struct{}

func (u user) setUserRouter(group *gin.RouterGroup) {
	group.GET("/login", u.login)
	group.GET("/signup", u.signup)
}

type loginRequest struct {
	Token int `form:"token" binding:"required"`
}

func (u user) login(c *gin.Context) {
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return
	}
	userId := req.Token
	user := models.User{}
	if err := user.Get(userId); err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, user)
}

type signupRequest struct {
	Name     string      `form:"name"  binding:"required"`
	Fullname null.String `form:"fullname"  binding:"-"`
}

func (u user) signup(c *gin.Context) {
	req := signupRequest{}
	if err := c.Bind(&req); err != nil {
		return
	}
	user := models.User{Name: req.Name, Fullname: req.Fullname}

	if err := user.Add(); err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, user)
}
