package customServiceError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUnknown = errors.New("unknown error. please, try again")
var ErrUserNotFound = errors.New("user not found")

var ErrTokenExpired = errors.New("token has expired")

var ErrTokenExpTimeBeforeIss = errors.New("token expiration time is before the issue time")
var ErrAccessTokenNotFound = errors.New("access token not found")
