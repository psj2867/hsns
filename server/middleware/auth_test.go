package middleware_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/psj2867/hsns/models"
	"github.com/psj2867/hsns/server/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	userid := "user"
	user := models.User{
		Name:   "asdf",
		UserId: userid,
	}
	token, err := middleware.GenerateJwtToken(middleware.UserInfoForToken{
		UserId: user.UserId,
	})
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("token: %v\n", token)
}

func TestVerifyToken(t *testing.T) {
	userid := "user"
	user := models.User{
		Name:   "asdf",
		UserId: userid,
	}
	token, err := middleware.GenerateJwtToken(middleware.UserInfoForToken{
		UserId: user.UserId,
	})
	if err != nil {
		t.Error(err)
		return
	}
	res, err := middleware.ParseJwt(token)
	if err != nil {
		t.Error(err)
		return
	}
	resUserId := res.Claims.(jwt.MapClaims)["userId"]
	assert.Equal(t, userid, resUserId)
}

func TestVerifyToken_time_fail(t *testing.T) {
	userid := "user"
	user := models.User{
		Name:   "asdf",
		UserId: userid,
	}
	token, err := middleware.GenerateJwtToken(middleware.UserInfoForToken{
		UserId: user.UserId,
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = middleware.ParseJwt(token, jwt.WithTimeFunc(func() time.Time {
		return time.Now().AddDate(0, 1, 0)
	}))
	assert.NotNil(t, err)
}
