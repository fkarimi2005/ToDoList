package errs

import "errors"

var (
	ErrUserAlreadyExists         = errors.New(`user already exists`)
	ErrNotFoundID                = errors.New("not found so ID")
	ErrAlreadyDeleted            = errors.New("this is task already deleted")
	ErrInvalidID                 = errors.New("invalid id")
	ErrValidationFailed          = errors.New("validation failed: Invalid input data")
	ErrSomethingWentWrong        = errors.New("something went wrong")
	ErrUserNotFound              = errors.New("user not found")
	ErrInvalidOperationType      = errors.New("invalid operation type")
	ErrNoPermissionsToCreateTask = errors.New("no permissions to create task this account")
	ErrTaskNotFound              = errors.New("task not found")
	ErrIncorrectLoginOrPassword  = errors.New("incorrect  login or password ")
	ErrNotFoud                   = errors.New("not foud")
	ErrNotAccess                 = errors.New("you are not authorized to access this resource")
)
