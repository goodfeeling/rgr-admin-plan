package errors

import (
	"errors"
	"net/http"
)

type ErrorType string
type ErrorMessage string

const (
	NotFound        ErrorType    = "NotFound"
	notFoundMessage ErrorMessage = "record not found"

	ValidationError        ErrorType    = "ValidationError"
	validationErrorMessage ErrorMessage = "validation error"

	ResourceAlreadyExists     ErrorType    = "ResourceAlreadyExists"
	alreadyExistsErrorMessage ErrorMessage = "resource already exists"

	RepositoryError        ErrorType    = "RepositoryError"
	repositoryErrorMessage ErrorMessage = "error in repository operation"

	NotAuthenticated             ErrorType    = "NotAuthenticated"
	notAuthenticatedErrorMessage ErrorMessage = "not Authenticated"

	TokenGeneratorError        ErrorType    = "TokenGeneratorError"
	tokenGeneratorErrorMessage ErrorMessage = "error in token generation"

	NotAuthorized             ErrorType    = "NotAuthorized"
	notAuthorizedErrorMessage ErrorMessage = "not authorized"

	UnknownError        ErrorType    = "UnknownError"
	unknownErrorMessage ErrorMessage = "something went wrong"

	UserNotExist ErrorType = "The user does not exist"
	UserExists   ErrorType = "The user already exists"

	TokenError        ErrorType    = "TokenError"
	TokenErrorMessage ErrorMessage = "error in token"

	TokenExpired        ErrorType    = "TokenIsExpired"
	TokenExpiredMessage ErrorMessage = "token is expired"

	UploadError        ErrorType    = "UploadError"
	UploadErrorMessage ErrorMessage = "upload has error"

	CaptchaError        ErrorType    = "CaptchaError"
	CaptchaErrorMessage ErrorMessage = "captcha error"
)

type AppError struct {
	Err  error
	Type ErrorType
}

func NewAppError(err error, errType ErrorType) *AppError {
	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func NewAppErrorWithType(errType ErrorType) *AppError {
	var err error

	switch errType {
	case NotFound:
		err = errors.New(string(notFoundMessage))
	case ValidationError:
		err = errors.New(string(validationErrorMessage))
	case ResourceAlreadyExists:
		err = errors.New(string(alreadyExistsErrorMessage))
	case RepositoryError:
		err = errors.New(string(repositoryErrorMessage))
	case NotAuthenticated:
		err = errors.New(string(notAuthenticatedErrorMessage))
	case NotAuthorized:
		err = errors.New(string(notAuthorizedErrorMessage))
	case TokenGeneratorError:
		err = errors.New(string(tokenGeneratorErrorMessage))
	case TokenError:
		err = errors.New(string(TokenErrorMessage))
	case TokenExpired:
		err = errors.New(string(TokenExpiredMessage))
	case UploadError:
		err = errors.New(string(UploadErrorMessage))
	case CaptchaError:
		err = errors.New(string(CaptchaErrorMessage))
	default:
		err = errors.New(string(unknownErrorMessage))
	}

	return &AppError{
		Err:  err,
		Type: errType,
	}
}

func (appErr *AppError) Error() string {
	return appErr.Err.Error()
}

// AppErrorToHTTP maps an AppError to an HTTP status code and message
func AppErrorToHTTP(appErr *AppError) (int, string) {
	switch appErr.Type {
	case NotFound:
		return http.StatusNotFound, appErr.Error()
	case ValidationError:
		return http.StatusBadRequest, appErr.Error()
	case RepositoryError:
		return http.StatusInternalServerError, appErr.Error()
	case NotAuthenticated:
		return http.StatusUnauthorized, appErr.Error()
	case NotAuthorized:
		return http.StatusForbidden, appErr.Error()
	case TokenError:
		return http.StatusUnauthorized, appErr.Error()
	case TokenExpired:
		return http.StatusUnauthorized, appErr.Error()
	case UploadError:
		return http.StatusBadRequest, appErr.Error()
	case CaptchaError:
		return http.StatusBadRequest, appErr.Error()
	default:
		return http.StatusInternalServerError, "Internal Server Error"
	}
}
