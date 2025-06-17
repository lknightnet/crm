package customRepositoryError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUserNotFound = errors.New("user not found")
var ErrUsersNotFound = errors.New("users not found")
