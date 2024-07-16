package router_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/psj2867/hsns/config"
	"github.com/psj2867/hsns/server/middleware"
	"github.com/psj2867/hsns/util"
	"github.com/stretchr/testify/assert"
)

func TestRouteContentUpload_AuthFail(t *testing.T) {
	res, req := util.HttptestPost("/contents/upload", map[string]string{
		"name": "asdf",
	})
	s.ServeHTTP(res, req)
	assert.Equal(t, 401, res.Result().StatusCode)
}

func TestRouteContentUpload(t *testing.T) {
	loginUser, err := routeLogin(testUser)
	if !assert.NoError(t, err) {
		return
	}
	res, req := util.HttptestPost("/contents/upload", map[string]string{
		"content": "asdf",
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
	fmt.Printf("res.Body.String(): %v\n", res.Body.String())
}

func TestRouteContentUploadSuccess(t *testing.T) {
	// login
	loginUser, err := routeLogin(testUser)
	if !assert.NoError(t, err) {
		return
	}
	// upload request
	res, req := util.HttptestPost("/contents/upload", map[string]string{
		"content": "asdf",
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
	// mock of upload image server
	returnToken := mockImageUpload(res.Body.String())
	// submit image upload success
	res, req = util.HttptestPost("/contents/upload/success", map[string]string{
		"id":    "1",
		"token": returnToken,
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())

}
func mockImageUpload(res string) string {
	decoded, _ := config.UploadTokenEnDecoder{}.Decode([]byte(res))
	fmt.Printf("decoded: %v\n", decoded)
	docodedRes := map[string]string{}
	json.Unmarshal(decoded, &docodedRes)
	fmt.Printf("docodedRes: %v\n", docodedRes)

	return ""
}
