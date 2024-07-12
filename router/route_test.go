package router_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/server"
	"github.com/psj2867/hsns/server/middleware"
	"github.com/psj2867/hsns/util"
)

var s *gin.Engine

func TestMain(m *testing.M) {
	config.SetTestDb()
	s = server.InitServer()
	gin.SetMode("test")
	defer server.DeferServer(s)
	os.Exit(m.Run())
}

func TestRoute(t *testing.T) {
	w, req := util.HttptestGet("/", nil)
	s.ServeHTTP(w, req)
	fmt.Printf("w: %v\n", w)
	// assert.Equal(t, 1, 1)
}

func routeUserSignup(userId string) *httptest.ResponseRecorder {
	w, req := util.HttptestPost("/user/signup", map[string]string{
		"name":   "asdf",
		"userid": userId,
	})
	s.ServeHTTP(w, req)
	return w
}
func TestRouteUserSignup(t *testing.T) {
	w := routeUserSignup("TestRouteUserSignup")
	assert.Equal(t, 200, w.Result().StatusCode, w.Body.String())
}
func TestRouteUserLogin(t *testing.T) {
	routeUserSignup("TestRouteUserLogin")
	w, req := util.HttptestGet("/user/login", map[string]string{
		"id": "TestRouteUserLogin",
	})
	s.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Result().StatusCode, w.Body.String())
}

func TestRouteUserAuthMiddleware(t *testing.T) {
	w := routeUserSignup("TestRouteUserAuthMiddleware")
	if !assert.Equal(t, 200, w.Result().StatusCode, w.Body.String()) {
		return
	}
	w, req := util.HttptestGet("/user/login", map[string]string{
		"id": "TestRouteUserAuthMiddleware",
	})
	s.ServeHTTP(w, req)
	if !assert.Equal(t, 200, w.Result().StatusCode, w.Body.String()) {
		return
	}

	res, err := toJson(w)
	if err != nil {
		t.Error(err)
		return
	}
	token := res["token"].(string)
	user, err := routeUserMe(token)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("user: %v\n", user)
	assert.Equal(t, "TestRouteUserAuthMiddleware", user["userId"])
}
func toJson(response *httptest.ResponseRecorder) (map[string]any, error) {
	var result map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &result); err != nil {
		return nil, err
	}
	return result, nil
}
func routeUserMe(token string) (map[string]any, error) {
	w, req := util.HttptestGet("/user/me", nil)
	req.Header.Add(middleware.Auth_Header, token)
	s.ServeHTTP(w, req)
	if w.Result().StatusCode != 200 {
		return nil, fmt.Errorf("return code %d, body: %s", w.Result().StatusCode, w.Body.String())
	}
	res, err := toJson(w)
	return res, err
}
