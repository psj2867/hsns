package router_test

import (
	"testing"

	"github.com/psj2867/hsns/server/middleware"
	"github.com/psj2867/hsns/util"
	"github.com/stretchr/testify/assert"
)

func TestFeed(t *testing.T) {
	loginUser, _ := routeLogin("qwer")

	res, req := util.HttptestGet("/contents/feed/main", nil)
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
}
