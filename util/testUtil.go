package util

import (
	"net/http"
	"net/http/httptest"
	"strings"
)

func HttptestPost(url string, values map[string]string) (w *httptest.ResponseRecorder, req *http.Request) {
	w = httptest.NewRecorder()
	var valueStrings []string
	for k, v := range values {
		valueStrings = append(valueStrings, k+"="+v)
	}
	body := strings.NewReader(strings.Join(valueStrings, "&"))
	req, _ = http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return
}
func HttptestGet(url string, values map[string]string) (w *httptest.ResponseRecorder, req *http.Request) {
	w = httptest.NewRecorder()
	var valueStrings []string
	for k, v := range values {
		valueStrings = append(valueStrings, k+"="+v)
	}
	body := strings.Join(valueStrings, "&")
	if len(body) > 0 {
		url = url + "?" + body
	}
	req, _ = http.NewRequest("GET", url, nil)
	return
}
