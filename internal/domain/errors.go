package domain

import "net/http"

var (
	ErrInvalidArgument = NewDomainError(10010, http.StatusUnprocessableEntity, "invalid argument")
	ErrNotFound        = NewDomainError(10000, http.StatusNotFound, "resource not found")
)
