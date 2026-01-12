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
	Key        string `json:"error_key"` //Hỗ trợ đa ngôn ngữ
}

// AppError(AppError(.....(error)))
func (err *AppError) RootError() error {
	if err, ok := err.RootErr.(*AppError); ok {
		return err.RootError()
	}
	return err.RootErr
}

// implement interface error
func (err *AppError) Error() string {
	return err.RootError().Error()
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

func NewCustomError(root error, msg, key string) *AppError {
	if root != nil {
		return NewErrorResponse(root, msg, root.Error(), key)
	}
	return NewErrorResponse(errors.New(msg), msg, msg, key)
}

func ErrDB(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong with DB", err.Error(), "DbError")
}

func ErrInvalidRequest(err error) *AppError {
	return NewErrorResponse(err, "invalid request", err.Error(), "InvalidRequest")
}

func ErrInternal(err error) *AppError {
	return NewFullErrorResponse(http.StatusInternalServerError, err, "something went wrong in the server", err.Error(), "InternalError")
}

func ErrCannotListEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot list %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotList%s", entity))
}

func ErrCannotDeleteEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot delete %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotDelete%s", entity))
}

func ErrCannotUpdateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot update %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotUpdate%s", entity))
}

func ErrCannotGetEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot get %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotGet%s", entity))
}

func ErrEntityDeleted(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s deteted", strings.ToLower(entity)),
		fmt.Sprintf("Err%sDelete", entity))
}

func ErrEntityExisted(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s aleardy exists", strings.ToLower(entity)),
		fmt.Sprintf("Err%sAlreadyExist", entity))
}

func ErrEntityNotFound(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("%s not found", strings.ToLower(entity)),
		fmt.Sprintf("Err%sNotFound", entity))
}

func ErrCannotCreateEntity(entity string, err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("Cannot create %s", strings.ToLower(entity)),
		fmt.Sprintf("ErrCannotCreate%s", entity))
}

func ErrNoPermission(err error) *AppError {
	return NewCustomError(err, fmt.Sprintf("You have no permission"),
		fmt.Sprintf("ErrNoPermission"))
}

var RecordNotFound = errors.New("record not found")
