package httpx

import (
	"net/http"
	"testing"

	"github.com/yanmai/api/pkg/errs"
)

func TestStatusOf(t *testing.T) {
	cases := []struct {
		code errs.Code
		want int
	}{
		{errs.CodeOK, http.StatusOK},
		{errs.CodeInvalidParam, http.StatusBadRequest},
		{errs.CodeUnauthorized, http.StatusUnauthorized},
		{errs.CodeTokenMissing, http.StatusUnauthorized},
		{errs.CodeTokenInvalid, http.StatusUnauthorized},
		{errs.CodeTokenExpired, http.StatusUnauthorized},
		{errs.CodeTokenTypeMismatch, http.StatusUnauthorized},
		{errs.CodeCredentialInvalid, http.StatusUnauthorized},
		{errs.CodeForbidden, http.StatusForbidden},
		{errs.CodePermissionDenied, http.StatusForbidden},
		{errs.CodeUserDisabled, http.StatusForbidden},
		{errs.CodeNotFound, http.StatusNotFound},
		{errs.CodeUserNotFound, http.StatusNotFound},
		{errs.CodeConflict, http.StatusConflict},
		{errs.CodeUserExists, http.StatusConflict},
		{errs.CodeTooManyRequest, http.StatusTooManyRequests},
		{errs.CodeNotImplemented, http.StatusNotImplemented},
		{errs.CodeMethodNotAllow, http.StatusMethodNotAllowed},
		{errs.CodeInternal, http.StatusInternalServerError},
		{errs.CodeUnknown, http.StatusInternalServerError},
	}

	for _, c := range cases {
		if got := StatusOf(c.code); got != c.want {
			t.Errorf("StatusOf(%d) = %d, 期望 %d", c.code, got, c.want)
		}
	}
}
