package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notificationService *services.NotificationService
}

func NewNotificationHandler(notificationService *services.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

func (h *NotificationHandler) NotificationsPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/notifications.html", gin.H{
		"title": "Notifications",
	})
}

func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	var request dtos.NotificationSearchRequest
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&request)) {
		return
	}

	response, err := h.notificationService.GetUserNotifications(c.Request.Context(), userID, request.Limit, request.Offset)
	if err != nil {
		appErrors.RespondCustomError(c, err, "Failed to fetch notifications")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}
	notificationIDParam := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.notificationService.MarkAsRead(c.Request.Context(), userID, uint(notificationID)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to mark notification as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	if err := h.notificationService.MarkAllAsRead(c.Request.Context(), userID); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to mark all notifications as read")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "All notifications marked as read"})
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}
	notificationIDParam := c.Param("id")
	notificationID, err := strconv.ParseUint(notificationIDParam, 10, 32)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid notification ID")
		return
	}

	if err := h.notificationService.DeleteNotification(c.Request.Context(), userID, uint(notificationID)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to delete notification")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}
