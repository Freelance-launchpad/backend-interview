package jerror

import (
	"errors"
	"fmt"
	"net/http"
)

type Details struct {
	Code  string `json:"code"`
	Field string `json:"field"`
	Err   error  `json:"-"`
}

func (d Details) String() string {
	var errStr string
	if d.Err != nil {
		errStr = ": " + d.Err.Error()
	}
	return d.Code + "(" + d.Field + ")" + errStr
}

type Error struct {
	errType    error
	Err        error
	Details    []Details
	HTTPStatus int
}

//nolint:staticcheck // ST1012: we want to keep the previous variable name from jerror/v1 and don't change the existing codebase
var (
	NotFoundError           = fmt.Errorf("jerror_not_found")
	PermanentlyDeletedError = fmt.Errorf("jerror_permanently_deleted")
	AlreadyExistsError      = fmt.Errorf("jerror_already_exists")
	ConflictError           = fmt.Errorf("jerror_conflict_error")
	ValidationError         = fmt.Errorf("jerror_validation")
	UnavailableError        = fmt.Errorf("jerror_unavailable")
	UnauthorizedError       = fmt.Errorf("jerror_unauthorized")
	PermissionDeniedError   = fmt.Errorf("jerror_permission_denied")
	BusinessError           = fmt.Errorf("jerror_business")
)

func (be *Error) Is(err error) bool {
	return errors.Is(be.Err, err) || errors.Is(be.errType, err)
}

func (be *Error) Error() string {
	return be.Err.Error()
}

func (be *Error) Unwrap() error {
	return be.Err
}

func NewNotFoundError(code string) error {
	return &Error{
		errType:    NotFoundError,
		HTTPStatus: http.StatusNotFound,
		Err:        errors.New(code),
	}
}

func NewPermanentlyDeletedError(code string) error {
	return &Error{
		errType:    PermanentlyDeletedError,
		HTTPStatus: http.StatusGone,
		Err:        errors.New(code),
	}
}

func NewAlreadyExistsError(code string) error {
	return &Error{
		errType:    AlreadyExistsError,
		HTTPStatus: http.StatusConflict,
		Err:        errors.New(code),
	}
}

func NewConflictError(code string) error {
	return &Error{
		errType:    ConflictError,
		HTTPStatus: http.StatusConflict,
		Err:        errors.New(code),
	}
}

func NewValidationError(code string, details ...Details) error {
	return &Error{
		errType:    ValidationError,
		HTTPStatus: http.StatusBadRequest,
		Err:        errors.New(code),
		Details:    details,
	}
}

func NewUnavailableError(code string) error {
	return &Error{
		errType:    UnavailableError,
		HTTPStatus: http.StatusServiceUnavailable,
		Err:        errors.New(code),
	}
}

// UnauthorizedError is currently handled as 400 and not 401 by [jgin.HandleError] because of a frontend limitation.
func NewUnauthorizedError(code string) error {
	return &Error{
		errType:    UnauthorizedError,
		HTTPStatus: http.StatusBadRequest,
		Err:        errors.New(code),
	}
}

func NewPermissionDeniedError(code string) error {
	return &Error{
		errType:    PermissionDeniedError,
		HTTPStatus: http.StatusForbidden,
		Err:        errors.New(code),
	}
}

func NewBusinessError(code string, details ...Details) error {
	return &Error{
		errType:    BusinessError,
		HTTPStatus: http.StatusUnprocessableEntity,
		Err:        errors.New(code),
		Details:    details,
	}
}
