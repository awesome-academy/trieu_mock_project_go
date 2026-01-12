package handlers

import (
	"log"
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"
	"trieu_mock_project_go/internal/utils"
	"trieu_mock_project_go/internal/websocket"

	"github.com/gin-gonic/gin"
	gorillaWebsocket "github.com/gorilla/websocket"
)

var upgrader = gorillaWebsocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now
	},
}

type NotificationHandler struct {
	notificationService *services.NotificationService
	hub                 *websocket.Hub
}

func NewNotificationHandler(notificationService *services.NotificationService, hub *websocket.Hub) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
		hub:                 hub,
	}
}

func (h *NotificationHandler) HandleWS(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		log.Println("WS: token missing")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	claims, err := utils.ParseJWTToken(token)
	if err != nil {
		log.Printf("WS: token invalid: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	conn, upgradeErr := upgrader.Upgrade(c.Writer, c.Request, nil)
	if upgradeErr != nil {
		log.Printf("WS: upgrade error: %v", upgradeErr)
		return
	}

	client := &websocket.Client{
		UserID: claims.UserID,
		Conn:   conn,
		Send:   make(chan *websocket.NotificationMessage, 256),
	}
	h.hub.Register(client)

	go client.ReadPump(h.hub)
	go client.WritePump(h.hub)
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

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		appErrors.RespondError(c, http.StatusUnauthorized, "Unauthorized access")
		return
	}

	count, err := h.notificationService.GetUnreadCount(c.Request.Context(), userID)
	if err != nil {
		appErrors.RespondCustomError(c, err, "Failed to fetch unread count")
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
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
