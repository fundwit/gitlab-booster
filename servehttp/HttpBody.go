package servehttp

import "github.com/fundwit/go-commons/types"

// ErrorBody 错误信息响应结构
type ErrorBody struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// SuccessBody 当个资源数据响应结构
type SuccessBody struct {
	Data interface{} `json:"data"`
}

// SuccessPagedBody 资源列表数据响应结构
type SuccessPagedBody struct {
	Data  interface{} `json:"data"`
	Total types.ID    `json:"total"`
}
