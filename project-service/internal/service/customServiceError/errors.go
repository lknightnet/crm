package customServiceError

import "errors"

var ErrUnknownError = errors.New("unknown error")
var ErrProjectNotFound = errors.New("project not found")
var ErrProjectsNotFound = errors.New("projects not found")
var ErrTaskNotFound = errors.New("task not found")
var ErrPermissionDenied = errors.New("permission denied")
var ErrTokenNotFound = errors.New("token not found")
var ErrTokenExpired = errors.New("token has expired")
var ErrAnotherTimerRunning = errors.New("another timer is running")
var ErrTimerNotFound = errors.New("timer not found")

var ErrInformationListNotFound = errors.New("lists not found")
