// Package errs 定义平台统一的业务错误码与错误类型。
//
// 设计约定：
//   - 错误码为 5 位整数，按业务域分段：1xxxx 通用 / 2xxxx 账号认证 / 3xxxx 课程 …
//   - 业务错误一律用 *errs.Error 表达，禁止把裸 error 直接抛给 HTTP 层
//   - 业务层不感知 HTTP，Code → HTTP Status 的映射集中在 pkg/httpx
package errs

import (
	"errors"
	"fmt"
)

// Code 业务错误码。
type Code int

// 通用错误码（1xxxx）。
const (
	CodeOK Code = 0

	CodeUnknown        Code = 10000
	CodeInvalidParam   Code = 10001
	CodeNotFound       Code = 10002
	CodeUnauthorized   Code = 10003
	CodeForbidden      Code = 10004
	CodeTooManyRequest Code = 10005
	CodeInternal       Code = 10006
	CodeConflict       Code = 10007
	CodeNotImplemented Code = 10008
	CodeMethodNotAllow Code = 10009
)

// 账号与认证错误码（2xxxx）。
const (
	CodeUserNotFound      Code = 20001
	CodeUserExists        Code = 20002
	CodeCredentialInvalid Code = 20003
	CodeUserDisabled      Code = 20004
	CodeTokenMissing      Code = 20005
	CodeTokenInvalid      Code = 20006
	CodeTokenExpired      Code = 20007
	CodeTokenTypeMismatch Code = 20008
	CodePermissionDenied  Code = 20009
)

var text = map[Code]string{
	CodeOK: "成功",

	CodeUnknown:        "服务开小差了，请稍后重试",
	CodeInvalidParam:   "请求参数不合法",
	CodeNotFound:       "请求的资源不存在",
	CodeUnauthorized:   "登录状态已失效，请重新登录",
	CodeForbidden:      "没有权限执行该操作",
	CodeTooManyRequest: "请求过于频繁，请稍后重试",
	CodeInternal:       "服务内部错误",
	CodeConflict:       "资源状态冲突",
	CodeNotImplemented: "功能尚未实现",
	CodeMethodNotAllow: "请求方法不被支持",

	CodeUserNotFound:      "账号不存在",
	CodeUserExists:        "账号已被注册",
	CodeCredentialInvalid: "账号或密码不正确",
	CodeUserDisabled:      "账号已被禁用，请联系管理员",
	CodeTokenMissing:      "缺少访问凭证",
	CodeTokenInvalid:      "访问凭证无效",
	CodeTokenExpired:      "访问凭证已过期",
	CodeTokenTypeMismatch: "访问凭证类型不正确",
	CodePermissionDenied:  "角色权限不足",
}

// Text 返回错误码的默认文案。
func (c Code) Text() string {
	if s, ok := text[c]; ok {
		return s
	}
	return text[CodeUnknown]
}

// Error 平台统一业务错误。
type Error struct {
	Code    Code
	Message string
	cause   error
}

// Error 实现 error 接口。
func (e *Error) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

// Unwrap 支持 errors.Is / errors.As 向上追溯。
func (e *Error) Unwrap() error { return e.cause }

// Is 判断目标错误是否携带指定错误码。
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e.Code == t.Code
	}
	return false
}

// New 构造一个业务错误，文案取错误码默认值。
func New(code Code) *Error {
	return &Error{Code: code, Message: code.Text()}
}

// Newf 构造一个业务错误并自定义文案。
func Newf(code Code, format string, args ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, args...)}
}

// Wrap 把底层 error 包装成业务错误，保留调用链。
func Wrap(code Code, cause error) *Error {
	if cause == nil {
		return New(code)
	}
	return &Error{Code: code, Message: code.Text(), cause: cause}
}

// WithCause 追加上游错误。
func (e *Error) WithCause(cause error) *Error {
	clone := *e
	clone.cause = cause
	return &clone
}

// WithMessage 覆盖文案。
func (e *Error) WithMessage(format string, args ...any) *Error {
	clone := *e
	clone.Message = fmt.Sprintf(format, args...)
	return &clone
}

// From 把任意 error 归一化为 *Error：已是业务错误则原样返回，否则视为内部错误。
func From(err error) *Error {
	if err == nil {
		return nil
	}
	var e *Error
	if errors.As(err, &e) {
		return e
	}
	return Wrap(CodeInternal, err)
}

// Is 判断任意 error 是否匹配指定错误码（业务代码里做分支判断的推荐方式）。
func Is(err error, code Code) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Code == code
	}
	return false
}

// As 是 errors.As 的语法糖，便于业务层少写一次 import。
func As(err error, target **Error) bool { return errors.As(err, target) }
