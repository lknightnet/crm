package customRepositoryError

import "errors"

var ErrInvalidEmailOrPassword = errors.New("invalid email or password")
var ErrUserNotFound = errors.New("user not found")
var ErrRouteLinkNotFound = errors.New("routeLink not found")
