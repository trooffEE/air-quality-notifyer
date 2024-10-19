package exceptions

import "errors"

//var (
//	ErrInternalDBError = exceptions.New("Internal Server Error")
//	ErrUserNotFound        = exceptions.New("User was not Found")
//	ErrUserAlreadyExists   = exceptions.New("User Already Exists")
//	ErrFailedToBeCreated   = exceptions.New("User Failed to be Created")
//)

var (
	ErrInternalDBError = errors.New("Internal Server Error")
)
