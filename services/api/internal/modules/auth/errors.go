package auth

import "github.com/yanmai/api/pkg/errs"

// errRefreshForbidden 令牌本身有效，但账号已不可用（被禁用或删除）。
func errRefreshForbidden() error {
	return errs.New(errs.CodeUserDisabled)
}
