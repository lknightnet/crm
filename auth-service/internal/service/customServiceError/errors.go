package customServiceError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUnknown = errors.New("unknown error. please, try again")
var ErrUserNotFound = errors.New("user not found")
var ErrTokenNotValid = errors.New("token is not valid")
var ErrOneTimeTokenExpired = errors.New("token has expired")
var ErrOneTimeTokenNotFound = errors.New("token not found")
