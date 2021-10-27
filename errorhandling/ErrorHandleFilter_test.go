package errorhandling_test

import (
	"errors"
	"fmt"
	"gitlab-booster/errorhandling"
	"gitlab-booster/testinfra"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ErrorHandleFilter", func() {
	var router *gin.Engine
	BeforeEach(func() {
		router = gin.Default()
		router.Use(errorhandling.GolbalErrorHandlingFilter())
	})

	Context("panic处理", func() {
		It("应当能够处理panic error", func() {
			router.GET("/", func(c *gin.Context) { panic(fmt.Errorf("some error")) })
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusInternalServerError))
			Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "msg":"some error", "data":null}`))
		})

		It("应当能够处理panic 普通对象", func() {
			router.GET("/", func(c *gin.Context) { panic("some error") })
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusInternalServerError))
			Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "msg":"some error", "data":null}`))
		})

		It("应当能够处理panic 业务错误", func() {
			demoErr := &demoError{Message: "some message in demo error", Data: 1234}
			router.GET("/", func(c *gin.Context) { panic(demoErr) })
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(444))
			Expect(body).To(MatchJSON(`{"code":"common.demo", "msg":"demo error: some message in demo error", "data":1234}`))
		})

		It("应当无法处理panic nil", func() {
			router.GET("/", func(c *gin.Context) { panic(nil) })
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusOK))
			Expect(body).To(Equal(``))
		})
	})

	Context("gin.Error 处理", func() {
		It("应当能够处理gin.Context.Errors中的错误", func() {
			router.GET("/", func(c *gin.Context) {
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error2")})
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusInternalServerError))
			Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "msg":"error2", "data":null}`))
		})

		It("应该在panic 非null并且有gin.Error时，处理panic", func() {
			router.GET("/", func(c *gin.Context) {
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error2")})
				panic("panic error")
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusInternalServerError))
			Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "msg":"panic error", "data":null}`))
		})

		It("应该在panic nil并且有gin.Error时，处理当前的gin.Error", func() {
			router.GET("/", func(c *gin.Context) {
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error1")})
				c.Errors = append(c.Errors, &gin.Error{Err: errors.New("error2")})
				panic(nil)
			})
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			status, body := testinfra.ExecuteRequest(req, router)
			Expect(status).To(Equal(http.StatusInternalServerError))
			Expect(body).To(MatchJSON(`{"code":"common.internal_server_error", "msg":"error2", "data":null}`))
		})
	})

})

type demoError struct {
	Message string
	Data    interface{}
}

func (e *demoError) Error() string {
	return fmt.Sprintf("demo error: %s", e.Message)
}

func (e *demoError) Detail() *errorhandling.BizErrorDetail {
	return &errorhandling.BizErrorDetail{
		Status: 444, Code: "common.demo",
		Message: e.Error(), Data: e.Data, Cause: nil,
	}
}
