package customServiceError

import "errors"

var ErrUnknown = errors.New("unknown error. please, try again")
var ErrUserNotFound = errors.New("user not found")
var ErrTokenNotValid = errors.New("token is not valid")
var ErrTokenNotFound = errors.New("token not found")

var ErrTokenExpired = errors.New("token has expired")
var ErrUsersNotFound = errors.New("users not found")
