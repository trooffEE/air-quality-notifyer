package exceptions

import "errors"

var (
	ErrInternalDBError = errors.New("Internal Server Error")
)
