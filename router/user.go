package router

import (
	"github.com/gin-gonic/gin"
	"github.com/psj2867/hsns/models"
	"github.com/psj2867/hsns/server/middleware"
)

type userRouter struct{}

func (u userRouter) setRouter(group *gin.RouterGroup) {
	group.GET("/login", u.login)
	group.POST("/signup", u.signup)
	group.GET("/me", middleware.ShouldAuthUnauth(), u.me)
}

type loginRequest struct {
	UserId string `form:"id" binding:"required"`
}

func (u *userRouter) login(c *gin.Context) {
	req := loginRequest{}
	if err := c.Bind(&req); err != nil {
		return
	}
	userId := req.UserId
	user := models.User{}
	if err := user.GetByUserId(userId); err != nil {
		c.JSON(403, err.Error())
		return
	}
	token, err := generateToken(&user)
	if err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, map[string]interface{}{
		"token":  token,
		"userId": user.Id,
	})
}

func generateToken(user *models.User) (string, error) {
	return middleware.GenerateToken(middleware.UserInfoForToken{
		UserId: user.UserId,
	})
}

type signupRequest struct {
	Name   string `form:"name"  binding:"required"`
	UserId string `form:"userid"  binding:"required"`
}

func (u *userRouter) signup(c *gin.Context) {
	req := signupRequest{}
	if err := c.Bind(&req); err != nil {
		return
	}
	user := models.User{Name: req.Name, UserId: req.UserId}

	if err := user.Add(); err != nil {
		c.JSON(403, err.Error())
		return
	}
	c.JSON(200, user)
}

func (u *userRouter) me(c *gin.Context) {
	userId, _ := middleware.GetAuthInfoByKey(c, "userId")
	user := models.User{}
	user.GetByUserId(userId.(string))
	c.JSON(200, user)
}
