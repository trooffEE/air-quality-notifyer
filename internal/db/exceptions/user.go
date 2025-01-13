package exceptions

import "errors"

var (
	UserNotFound = errors.New("User not found")
)
