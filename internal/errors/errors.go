package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
)

// Errors definitions
var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrNotFound                = errors.New("not found")
	ErrForbidden               = errors.New("forbidden")
	ErrMissingAuthHeader       = errors.New("missing authorization header")
	ErrInvalidAuthHeader       = errors.New("invalid authorization header format")
	ErrInvalidToken            = errors.New("invalid or expired token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrEmailAlreadyExists      = errors.New("email already exists")
	ErrUserNotFound            = errors.New("user not found")
	ErrPositionAlreadyExists   = errors.New("position with name already exists")
	ErrPositionInUse           = errors.New("position is assigned to one or more users")
	ErrSkillAlreadyExists      = errors.New("skill with name already exists")
	ErrSkillInUse              = errors.New("skill is assigned to one or more users")
)

// Error response
type APIErrorResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

func RespondError(
	c *gin.Context,
	status int,
	message string,
	details ...interface{},
) {
	var detail interface{} = nil
	if len(details) > 0 {
		detail = details[0]
	}
	c.JSON(status, APIErrorResponse{
		Code:    status,
		Message: message,
		Details: detail,
	})
}

func HandleBindError(c *gin.Context, err error) bool {
	if err == nil {
		return false
	}

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		RespondError(
			c,
			http.StatusBadRequest,
			"Invalid request body format",
			nil,
		)
		return true
	}

	fields := make(map[string]string)
	for _, fieldErr := range validationErrs {
		field := fieldErr.Field()

		switch fieldErr.Tag() {
		case "required":
			fields[field] = "is required"
		case "email":
			fields[field] = "must be a valid email"
		case "min":
			fields[field] = "must be at least " + fieldErr.Param() + " characters"
		case "max":
			fields[field] = "must be at most " + fieldErr.Param() + " characters"
		default:
			fields[field] = "is invalid"
		}
	}

	RespondError(
		c,
		http.StatusBadRequest,
		"Validation failed",
		fields,
	)
	return true
}

func IsDuplicatedEntryError(err error) bool {
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		return true
	}
	return false
}
