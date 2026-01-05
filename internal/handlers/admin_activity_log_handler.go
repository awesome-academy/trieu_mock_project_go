package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type AdminActivityLogHandler struct {
	activityLogService *services.ActivityLogService
}

func NewAdminActivityLogHandler(activityLogService *services.ActivityLogService) *AdminActivityLogHandler {
	return &AdminActivityLogHandler{activityLogService: activityLogService}
}

func (h *AdminActivityLogHandler) AdminActivityLogsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_activity_logs.html", gin.H{
		"title":     "Admin Activity Logs Management",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminActivityLogHandler) ActivityLogsSearchPartial(c *gin.Context) {
	templateName := "partials/admin_activity_logs_search.html"
	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid query parameters")
		return
	}

	resp, err := h.activityLogService.SearchActivityLogs(c.Request.Context(), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load activity logs")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"logs": resp.Logs,
		"page": resp.Page,
	})
}

func (h *AdminActivityLogHandler) DeleteActivityLog(c *gin.Context) {
	activityLogIdParam := c.Param("activityLogId")
	activityLogId, err := strconv.Atoi(activityLogIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid activity log ID")
		return
	}

	if err := h.activityLogService.DeleteActivityLog(c.Request.Context(), uint(activityLogId)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to delete activity log")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Activity log deleted successfully"})
}
