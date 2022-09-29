package model

import "errors"

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = errors.New("internal server error")
	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = errors.New("your requested item is not found")
	// ErrConflict will throw if the current action already exists
	ErrConflict = errors.New("your item already exist")
	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = errors.New("given param is not valid")
	//ErrUnauthorized will throw if request lacks valid authentication credentials
	ErrUnauthorized = errors.New("unauthorized")
	//ErrServiceUnavailable will throw if server is not responding
	ErrServiceUnavailable = errors.New("service unavailable")
	//
	ErrForbidden = errors.New("forbidden Error")
	ErrInvalidCa = errors.New("invalid ca ceritficate")
)
