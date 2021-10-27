package upstream

import (
	"fmt"
	"gitlab-booster/errorhandling"
	"net/http"
)

// ErrUpstreamServiceConnectivity 依赖的上游服务可连接性问题
type ErrUpstreamServiceConnectivity struct {
	Service string
	Request string
	Err     error
}

func (e *ErrUpstreamServiceConnectivity) Error() string {
	return fmt.Sprintf("connectivity problem with upstream service '%s' when handling request '%s': %v", e.Service, e.Request, e.Err)
}
func (e *ErrUpstreamServiceConnectivity) Unwrap() error {
	return e.Err
}

// Detail 获取详细的业务错误数据
func (e *ErrUpstreamServiceConnectivity) Detail() *errorhandling.BizErrorDetail {
	return &errorhandling.BizErrorDetail{
		Status: http.StatusBadGateway, Code: "upstream.conneciton_failed",
		Message: e.Error(), Data: nil, Cause: e.Err,
	}
}
