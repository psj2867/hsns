package router_test

import (
	"testing"

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
	res, req := util.HttptestPost("/contents/upload", map[string]string{
		"name": "asdf",
	})
	s.ServeHTTP(res, req)
	assert.Equal(t, 401, res.Result().StatusCode)
}
