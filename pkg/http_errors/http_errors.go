package httpErrors

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	ErrBadRequest       = "Bad request"
	ErrAlreadyExists    = "Already exists"
	ErrNoSuchUser       = "User not found"
	ErrWrongCredentials = "Wrong Credentials"
	ErrNotFound         = "Not Found"
	ErrUnauthorized     = "Unauthorized"
	ErrForbidden        = "Forbidden"
	ErrBadQueryParams   = "Invalid query params"
	ErrRequestTimeout   = "Request Timeout"
	ErrInvalidEmail     = "Invalid email"
	ErrInvalidPassword  = "Invalid password"
	ErrInvalidField     = "Invalid field"
)

var (
	BadRequest            = errors.New("Bad request")
	WrongCredentials      = errors.New("Wrong Credentials")
	NotFound              = errors.New("Not Found")
	Unauthorized          = errors.New("Unauthorized")
	Forbidden             = errors.New("Forbidden")
	PermissionDenied      = errors.New("Permission Denied")
	ExpiredCSRFError      = errors.New("Expired CSRF token")
	WrongCSRFToken        = errors.New("Wrong CSRF token")
	CSRFNotPresented      = errors.New("CSRF not presented")
	NotRequiredFields     = errors.New("No such required fields")
	BadQueryParams        = errors.New("Invalid query params")
	InternalServerError   = errors.New("Internal Server Error")
	RequestTimeoutError   = errors.New("Request Timeout")
	ExistsEmailError      = errors.New("User with given email already exists")
	InvalidJWTToken       = errors.New("Invalid JWT token")
	InvalidJWTClaims      = errors.New("Invalid JWT claims")
	NotAllowedImageHeader = errors.New("Not allowed image header")
	NoCookie              = errors.New("not found cookie header")
	InvalidUUID           = errors.New("invalid uuid")
)

type RestErr interface{
	Status() int
	Error() string
	Causes() interface{}
	ErrBody() RestError
}

type RestError struct{
	ErrStatus	int			`json:"status,omitempty"`
	ErrError	string		`json:"error,omitempty"`
	ErrCauses	interface{}	`json:"err_causes,omitempty"`
}

func (e RestError) ErrBody() RestError {
	return e
}

func (e RestError) Error() string {
	return fmt.Sprintf("status: %d - errors: %s - causes: %v", e.ErrStatus, e.ErrError,e.ErrCauses)
}
func (e RestError) status() int {
	return e.ErrStatus
}

func (e RestError) Causes() interface{} {
	return e.ErrCauses
}

func NewRestError(status int, err string, causes interface{}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError: err,
		ErrCauses: causes,
	}
}

func NewRestErrorWithMessage(status int, err string, causes interface {}) RestErr {
	return RestError{
		ErrStatus: status,
		ErrError: err,
		ErrCauses: causes,
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr RestErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, errors.New("invalid json")
	}
	return apiErr, nil	
}

func NewBadRequestError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusBadRequest,
		ErrError:  BadRequest.Error(),
		ErrCauses: causes,
	}
}

func NewNotFoundError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusNotFound,
		ErrError:  NotFound.Error(),
		ErrCauses: causes,
	}
}

func NewUnauthorizedError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusUnauthorized,
		ErrError:  Unauthorized.Error(),
		ErrCauses: causes,
	}
}

func NewForbiddenError(causes interface{}) RestErr {
	return RestError{
		ErrStatus: http.StatusForbidden,
		ErrError:  Forbidden.Error(),
		ErrCauses: causes,
	}
}

func NewInternalServerError(causes interface{}) RestErr {
	result := RestError{
		ErrStatus: http.StatusInternalServerError,
		ErrError:  InternalServerError.Error(),
		ErrCauses: causes,
	}
	return result
}

func ParseErrors(err error) RestErr {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return NewRestError(http.StatusNotFound, ErrNotFound, nil)
	case errors.Is(err, context.DeadlineExceeded):
		return NewRestError(http.StatusRequestTimeout, ErrRequestTimeout, nil)
	case errors.Is(err, Unauthorized):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, nil)
	case errors.Is(err, WrongCredentials):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, nil)
	case strings.Contains(strings.ToLower(err.Error()), "sqlstate"):
		return parseSqlErrors(err)
	case strings.Contains(strings.ToLower(err.Error()), "field validation"):
		return parseValidatorError(err)
	case strings.Contains(strings.ToLower(err.Error()), "unmarshal"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err)
	case strings.Contains(strings.ToLower(err.Error()), "uuid"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, err)
	case strings.Contains(strings.ToLower(err.Error()), "cookie"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err)
	case strings.Contains(strings.ToLower(err.Error()), "token"):
		return NewRestError(http.StatusUnauthorized, ErrUnauthorized, err)
	case strings.Contains(strings.ToLower(err.Error()), "bcrypt"):
		return NewRestError(http.StatusBadRequest, ErrBadRequest, nil)
	default:
		if restErr, ok := err.(RestErr); ok {
			return restErr
		}
		return NewInternalServerError(err)
	}
}

func parseSqlErrors(err error) RestErr {
	return NewRestError(http.StatusBadRequest, ErrBadRequest, err)
}

func parseValidatorError(err error) RestErr {
	if strings.Contains(err.Error(), "Password") {
		return NewRestError(http.StatusBadRequest, ErrInvalidPassword, err)
	}

	if strings.Contains(err.Error(), "Email") {
		return NewRestError(http.StatusBadRequest, ErrInvalidEmail, err)
	}

	return NewRestError(http.StatusBadRequest, ErrInvalidField, err)
}

func ErrorResponse(err error) (int, interface{}) {
	return ParseErrors(err).Status(), ParseErrors(err)
}

func ErrorCtxResponse(ctx echo.Context, err error) error {
	restErr := ParseErrors(err)
	return ctx.JSON(restErr.Status(), restErr.ErrBody())
}