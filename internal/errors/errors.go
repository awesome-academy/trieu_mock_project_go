package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-sql-driver/mysql"
	csrf "github.com/utrack/gin-csrf"
)

type AppError struct {
	Status  int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(status int, message string) *AppError {
	return &AppError{Status: status, Message: message}
}

// Errors definitions
var (
	ErrInternalServerError                = NewAppError(http.StatusInternalServerError, "internal server error")
	ErrInvalidCredentials                 = NewAppError(http.StatusUnauthorized, "invalid credentials")
	ErrNotFound                           = NewAppError(http.StatusNotFound, "not found")
	ErrForbidden                          = NewAppError(http.StatusForbidden, "forbidden")
	ErrMissingAuthHeader                  = NewAppError(http.StatusUnauthorized, "missing authorization header")
	ErrInvalidAuthHeader                  = NewAppError(http.StatusUnauthorized, "invalid authorization header format")
	ErrInvalidToken                       = NewAppError(http.StatusUnauthorized, "invalid or expired token")
	ErrUnexpectedSigningMethod            = NewAppError(http.StatusUnauthorized, "unexpected signing method")
	ErrEmailAlreadyExists                 = NewAppError(http.StatusConflict, "email already exists")
	ErrUserNotFound                       = NewAppError(http.StatusNotFound, "user not found")
	ErrPositionNotFound                   = NewAppError(http.StatusNotFound, "position not found")
	ErrPositionAlreadyExists              = NewAppError(http.StatusConflict, "position with name already exists")
	ErrPositionInUse                      = NewAppError(http.StatusBadRequest, "position is assigned to one or more users")
	ErrSkillNotFound                      = NewAppError(http.StatusNotFound, "skill not found")
	ErrSkillAlreadyExists                 = NewAppError(http.StatusConflict, "skill with name already exists")
	ErrSkillInUse                         = NewAppError(http.StatusBadRequest, "skill is assigned to one or more users")
	ErrTeamAlreadyExists                  = NewAppError(http.StatusConflict, "team with name already exists")
	ErrTeamLeaderAlreadyInAnotherTeam     = NewAppError(http.StatusBadRequest, "team leader is already leading another team")
	ErrTeamNotFound                       = NewAppError(http.StatusNotFound, "team not found")
	ErrUserAlreadyInTeam                  = NewAppError(http.StatusBadRequest, "user is already a member of the team")
	ErrUserNotInTeam                      = NewAppError(http.StatusBadRequest, "user is not a member of the team")
	ErrCannotRemoveOrMoveTeamLeader       = NewAppError(http.StatusBadRequest, "cannot remove or move the team leader from the team")
	ErrCannotRemoveOrMoveProjectMember    = NewAppError(http.StatusBadRequest, "cannot remove or move the user because they are a member of a project in this team")
	ErrCannotDeleteUserBeingTeamLeader    = NewAppError(http.StatusBadRequest, "user cannot be deleted because they are a team leader")
	ErrProjectNotFound                    = NewAppError(http.StatusNotFound, "project not found")
	ErrProjectAlreadyExists               = NewAppError(http.StatusConflict, "project with name already exists")
	ErrCannotDeleteUserBeingProjectLeader = NewAppError(http.StatusBadRequest, "user cannot be deleted because they are a project leader")
	ErrCannotDeleteUserBeingProjectMember = NewAppError(http.StatusBadRequest, "user cannot be deleted because they are a project member")
	ErrActivityLogNotFound                = NewAppError(http.StatusNotFound, "activity log not found")
	ErrNoCSVDataToImport                  = NewAppError(http.StatusBadRequest, "no CSV data to import")
	ErrInvalidCSVFormat                   = NewAppError(http.StatusBadRequest, "invalid CSV format")
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

func RespondCustomError(c *gin.Context, err error, defaultMessage string) {
	if appErr, ok := err.(*AppError); ok {
		RespondError(c, appErr.Status, appErr.Message)
		return
	}

	RespondError(c, http.StatusInternalServerError, defaultMessage)
}

func RespondPageError(
	c *gin.Context,
	status int,
	templateName string,
	message string,
) {
	c.HTML(status, templateName, gin.H{
		"error": message,
	})
}

func RespondPageErrorWithCSRF(
	c *gin.Context,
	status int,
	templateName string,
	message string,
) {
	c.HTML(status, templateName, gin.H{
		"error":     message,
		"csrfToken": csrf.GetToken(c),
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
		case "ltfield":
			fields[field] = "must be less than " + fieldErr.Param()
		case "gtfield":
			fields[field] = "must be greater than " + fieldErr.Param()
		case "required_with_end_date":
			fields[field] = "is required when End Date is provided"
		case "gt_start_date":
			fields[field] = "must be greater than Start Date"
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
