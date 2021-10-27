package errorhandling

import (
	"errors"
	"fmt"
	"gitlab-booster/i18n"
	"gitlab-booster/utils"

	"github.com/gin-gonic/gin"
)

type errorBody struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// GolbalErrorHandlingFilter 全局错误处理
func GolbalErrorHandlingFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handleAll(c)
		c.Next() // execute all the handlers
	}
}

// HandleAll 支持处理 panic 和 c.Errors 的错误
func handleAll(c *gin.Context) {
	if err := recover(); err != nil {
		if e, ok := err.(error); ok {
			handleError(c, e)
		} else {
			c.JSON(500, &errorBody{Code: i18n.CommonInternalServerErrorCode, Msg: fmt.Sprintf("%s", err)})
		}
	} else {
		if err := c.Errors.Last(); err != nil {
			handleError(c, err)
		}
	}
}

// HandleError 处理错误信息
func handleError(c *gin.Context, err error) {
	// log out
	utils.Log.Error(err)

	genericErr := err
	var ginErr *gin.Error
	if errors.As(err, &ginErr) {
		genericErr = ginErr.Err
	}

	// write response
	if bizErr, ok := genericErr.(BizError); ok {
		detail := bizErr.Detail()
		c.JSON(detail.Status, &errorBody{Code: detail.Code, Msg: detail.Message, Data: detail.Data})
		return
	}

	c.JSON(500, &errorBody{Code: i18n.CommonInternalServerErrorCode, Msg: genericErr.Error()})
}
