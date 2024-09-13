package errApp

import "errors"

var (
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrUserNotFound        = errors.New("User was not Found")
	ErrUserAlreadyExists   = errors.New("User Already Exists")
	ErrFailedToBeCreated   = errors.New("User Failed to be Created")
)
