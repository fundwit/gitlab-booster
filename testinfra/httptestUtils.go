package testinfra

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// ExecuteRequest 执行请求对象，并返回响应的状态码和响应体字符串
func ExecuteRequest(req *http.Request, engine *gin.Engine) (int, string) {
	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	return resp.StatusCode, string(bodyBytes)
}
