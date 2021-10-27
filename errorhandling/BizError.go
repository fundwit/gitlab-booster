package errorhandling

// BizError 代表业务错误
type BizError interface {
	// Detail 获取详细的业务错误数据
	Detail() *BizErrorDetail
}

// BizErrorDetail 业务错误的详细信息
type BizErrorDetail struct {
	Status  int
	Code    string
	Message string

	Data  interface{}
	Cause error
}
