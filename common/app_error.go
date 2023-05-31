package common

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	RootErr    error  `json:"-"`
	Message    string `json:"message"`
	Log        string `json:"log"`
	Key        string `json:"error_key"`
}

func NewErrorResponse(root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewFullErrorResponse(statusCode int, root error, msg, log, key string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		RootErr:    root,
		Message:    msg,
		Log:        log,
		Key:        key,
	}
}

func NewUnauthorized(root error, msg, key string) *AppError {
	return &AppError{
		StatusCode: http.StatusUnauthorized,
		RootErr:    root,
		Message:    msg,
		Key:        key,
	}
}

func NewCustomError(root error, msg string, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}

	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func (e *AppError) RootError() error {
	if err, ok := e.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return e.RootErr
}
func (e *AppError) Error() string {
	return e.RootError().Error()
}

func ErrorDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "some thing went wrong with DB", err.Error(), "DB_ERROR")
}

func ErrIvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "ivalid request", err.Error(), "ErrIvalidRequest")
}
func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in the server", err.Error(), "ErrInternal")
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("cannot list %s", strings.ToLower(entity)),
		fmt.Sprintln("ErrCannotList%s", entity),
	)
}
func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintln("ErrCannotDelete%s", entity),
	)
}
func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("cannot update %s", strings.ToLower(entity)),
		fmt.Sprintln("ErrCannotUpdate%s", entity),
	)
}
func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("cannot get %s", strings.ToLower(entity)),
		fmt.Sprintln("ErrCannotGet%s", entity),
	)
}

func ErrEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("%s deleted", strings.ToLower(entity)),
		fmt.Sprintln("Err%sDeleted", entity),
	)
}
func ErrEntityExisted(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("%s already exists", strings.ToLower(entity)),
		fmt.Sprintln("Err%sAlreadyExists", entity),
	)
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("%s not found", strings.ToLower(entity)),
		fmt.Sprintln("Err%sNotFound", entity),
	)
}
func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("cannot create %s", strings.ToLower(entity)),
		fmt.Sprintln("ErrCannotCreate%s", entity),
	)
}
func ErrNoPermission(err error) *AppError {
	return NewCustomError(
		err,
		fmt.Sprintln("you have no permission"),
		fmt.Sprintln("ErrNoPermission"),
	)
}

var RecordNotFound = errors.New("record not found")
