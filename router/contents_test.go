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
		"images":  "3",
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
	// mock of upload image server
	uploadToken := res.Body.String()
	returnToken := mockImageUpload(uploadToken)
	// submit image upload success
	res, req = util.HttptestPost("/contents/upload/success", map[string]string{
		"token": returnToken,
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
}
func mockImageUpload(res string) string {
	decoded, err := config.UploadTokenEnDecoder.Decode([]byte(res))
	if err != nil {
		panic(err)
	}
	fmt.Printf("decoded: %v\n", string(decoded))
	uploadTokenData := map[string]any{}
	if err := json.Unmarshal(decoded, &uploadTokenData); err != nil {
		panic(err)
	}
	fmt.Printf("docodedRes: %v\n", uploadTokenData)
	jsonB, _ := json.Marshal(map[string]any{
		"uuid":           uploadTokenData["uuid"],
		"requestImages":  uploadTokenData["images"],
		"uploadedImages": uploadTokenData["images"],
	})
	r, _ := config.ReturnTokenEnDecoder.Encode(jsonB)
	return string(r)
}

func TestRouteContentUploadFail(t *testing.T) {
	// login
	loginUser, err := routeLogin(testUser)
	if !assert.NoError(t, err) {
		return
	}
	// upload request
	res, req := util.HttptestPost("/contents/upload", map[string]string{
		"content": "asdf",
		"images":  "3",
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
	// mock of upload image server
	uploadToken := res.Body.String()
	// submit image upload success
	res, req = util.HttptestPost("/contents/upload/success", map[string]string{
		"token": uploadToken,
	})
	req.Header.Add(middleware.Auth_Header, loginUser["token"].(string))
	s.ServeHTTP(res, req)
	assert.Equal(t, 200, res.Result().StatusCode, res.Body.String())
}
