// Package httpx 是 Gin 与业务错误之间的适配层。
//
// 职责边界：
//   - 把 *errs.Error 翻译成 HTTP 状态码 + 统一响应体
//   - 提供参数绑定（含校验失败归一化）与 traceId 读取
//
// 不含任何业务逻辑。
package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/yanmai/api/pkg/errs"
	"github.com/yanmai/api/pkg/response"
)

// TraceIDKey 上下文与响应头中 traceId 的键名。
const (
	TraceIDKey    = "traceId"
	TraceIDHeader = "X-Trace-Id"
)

// TraceID 从上下文读取当前请求的 traceId。
func TraceID(c *gin.Context) string {
	if v, ok := c.Get(TraceIDKey); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// OK 返回 200 + 成功响应体。
func OK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, response.New(int(errs.CodeOK), errs.CodeOK.Text(), data, TraceID(c)))
}

// Created 返回 201 + 成功响应体。
func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, response.New(int(errs.CodeOK), errs.CodeOK.Text(), data, TraceID(c)))
}

// NoContent 返回 204，用于无响应体的成功操作。
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Page 返回 200 + 分页响应体。
func Page(c *gin.Context, list any, total int64, page, pageSize int) {
	OK(c, response.NewPage(list, total, page, pageSize))
}

// Fail 返回业务错误对应的响应。
//
// 若 err 不是业务错误，会被归一化为 CodeInternal，且原始信息只落到日志（由 middleware 处理），
// 对外统一提示「服务内部错误」，避免泄露实现细节。
func Fail(c *gin.Context, err error) {
	e := errs.From(err)
	if e == nil {
		e = errs.New(errs.CodeInternal)
	}
	c.AbortWithStatusJSON(StatusOf(e.Code), response.New(int(e.Code), e.Message, nil, TraceID(c)))
}

// StatusOf 把业务错误码映射为 HTTP 状态码。
func StatusOf(code errs.Code) int {
	switch code {
	case errs.CodeOK:
		return http.StatusOK
	case errs.CodeInvalidParam:
		return http.StatusBadRequest
	case errs.CodeUnauthorized,
		errs.CodeTokenMissing,
		errs.CodeTokenInvalid,
		errs.CodeTokenExpired,
		errs.CodeTokenTypeMismatch,
		errs.CodeCredentialInvalid:
		return http.StatusUnauthorized
	case errs.CodeForbidden,
		errs.CodePermissionDenied,
		errs.CodeUserDisabled:
		return http.StatusForbidden
	case errs.CodeNotFound, errs.CodeUserNotFound:
		return http.StatusNotFound
	case errs.CodeConflict, errs.CodeUserExists:
		return http.StatusConflict
	case errs.CodeTooManyRequest:
		return http.StatusTooManyRequests
	case errs.CodeNotImplemented:
		return http.StatusNotImplemented
	case errs.CodeMethodNotAllow:
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

// Bind 绑定并校验请求体；失败时直接写入 400 响应并返回 false。
//
// 校验错误会被归一化为「字段名: 规则」，便于前端定位。
func Bind(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		Fail(c, errs.Wrap(errs.CodeInvalidParam, err).WithMessage("%s", FormatBindError(err)))
		return false
	}
	return true
}

// BindQuery 绑定并校验查询参数；失败时直接写入 400 响应并返回 false。
func BindQuery(c *gin.Context, obj any) bool {
	if err := c.ShouldBindQuery(obj); err != nil {
		Fail(c, errs.Wrap(errs.CodeInvalidParam, err).WithMessage("%s", FormatBindError(err)))
		return false
	}
	return true
}

// FormatBindError 把 binding 错误整理成人话。
func FormatBindError(err error) string {
	var ve interface{ Error() string }
	if ok := asValidatorErrors(err, &ve); ok {
		return ve.Error()
	}
	return errs.CodeInvalidParam.Text()
}
