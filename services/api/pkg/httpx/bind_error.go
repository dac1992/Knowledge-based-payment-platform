package httpx

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// asValidatorErrors 尝试把错误还原为可读的字段校验信息。
func asValidatorErrors(err error, out *interface{ Error() string }) bool {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return false
	}
	parts := make([]string, 0, len(verrs))
	for _, fe := range verrs {
		parts = append(parts, fmt.Sprintf("%s 不满足 %s", lowerFirst(fe.Field()), describeTag(fe)))
	}
	msg := strings.Join(parts, "；")
	*out = textError(msg)
	return true
}

type textError string

func (t textError) Error() string { return string(t) }

func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func describeTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "必填"
	case "email":
		return "邮箱格式"
	case "min":
		return "最小长度/值 " + fe.Param()
	case "max":
		return "最大长度/值 " + fe.Param()
	case "len":
		return "长度必须为 " + fe.Param()
	case "oneof":
		return "取值必须是 [" + fe.Param() + "] 之一"
	case "e164":
		return "E.164 手机号格式"
	case "url":
		return "URL 格式"
	case "gte":
		return "大于等于 " + fe.Param()
	case "lte":
		return "小于等于 " + fe.Param()
	default:
		return fe.Tag()
	}
}
