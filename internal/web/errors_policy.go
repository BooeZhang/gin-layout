package web

import (
	"net/http"

	"gin-layout/internal/policy"
)

var policyErrorMappings = []errorMapping{
	newErrorMapping(policy.ErrPermissionDenied, http.StatusForbidden, 30130, "没有权限"),
}
