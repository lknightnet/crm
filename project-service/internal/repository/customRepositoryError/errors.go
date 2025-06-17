package customRepositoryError

import "errors"

var ErrProjectNotFound = errors.New("project not found")
var ErrProjectsNotFound = errors.New("projects not found")
var ErrTaskNotFound = errors.New("task not found")
var ErrTimerNotFound = errors.New("timer not found")

var ErrInformationListNotFound = errors.New("lists not found")
