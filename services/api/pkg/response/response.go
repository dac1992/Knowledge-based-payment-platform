// Package response 定义平台统一的 HTTP 响应信封。
//
// 所有接口（含错误）统一返回：
//
//	{ "code": 0, "message": "成功", "data": {...}, "traceId": "..." }
//
// 其中 code 为 errs.Code，0 表示成功；HTTP 状态码另行表达语义（见 pkg/httpx）。
package response

// Body 统一响应体。
type Body struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"traceId,omitempty"`
}

// New 构造响应体。
func New(code int, message string, data any, traceID string) Body {
	return Body{Code: code, Message: message, Data: data, TraceID: traceID}
}

// Page 分页数据容器。
type Page struct {
	List     any   `json:"list"`
	Total    int64 `json:"total"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

// NewPage 构造分页数据；page/pageSize 会被规范化为合法值。
func NewPage(list any, total int64, page, pageSize int) Page {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return Page{List: list, Total: total, Page: page, PageSize: pageSize}
}
