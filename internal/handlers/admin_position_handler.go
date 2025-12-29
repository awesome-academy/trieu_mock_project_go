package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type AdminPositionHandler struct {
	positionsService *services.PositionService
}

func NewAdminPositionHandler(positionService *services.PositionService) *AdminPositionHandler {
	return &AdminPositionHandler{positionsService: positionService}
}

func (h *AdminPositionHandler) ListPositionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_positions.html", gin.H{
		"title": "Admin Positions Management",
	})
}

func (h *AdminPositionHandler) PositionSearchPartial(c *gin.Context) {
	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		c.HTML(http.StatusBadRequest, "partials/admin_positions_search.html", gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	resp, err := h.positionsService.SearchPositions(c.Request.Context(), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "partials/admin_positions_search.html", gin.H{
			"error": "Failed to load positions",
		})
		return
	}

	c.HTML(http.StatusOK, "partials/admin_positions_search.html", gin.H{
		"positions": resp.Positions,
		"page":      resp.Page,
	})
}

func (h *AdminPositionHandler) CreatePositionPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_position_create.html", gin.H{
		"title": "Create Position",
	})
}

func (h *AdminPositionHandler) CreatePosition(c *gin.Context) {
	var request dtos.CreateOrUpdatePositionRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.positionsService.CreatePosition(c.Request.Context(), request); err != nil {
		if err == appErrors.ErrPositionAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to create position")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position created successfully"})
}

func (h *AdminPositionHandler) EditPositionPage(c *gin.Context) {
	positionIdParam := c.Param("positionId")
	positionId, err := strconv.Atoi(positionIdParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "pages/admin_position_edit.html", gin.H{
			"title": "Edit Position",
			"error": "Invalid position ID",
		})
		return
	}

	position, err := h.positionsService.GetPositionByID(c.Request.Context(), uint(positionId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_position_edit.html", gin.H{
			"title": "Edit Position",
			"error": "Failed to load position details",
		})
		return
	}

	c.HTML(http.StatusOK, "pages/admin_position_edit.html", gin.H{
		"title":    "Edit Position",
		"position": position,
	})
}

func (h *AdminPositionHandler) UpdatePosition(c *gin.Context) {
	positionIdParam := c.Param("positionId")
	positionId, err := strconv.Atoi(positionIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid position ID")
		return
	}

	var request dtos.CreateOrUpdatePositionRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.positionsService.UpdatePosition(c.Request.Context(), uint(positionId), request); err != nil {
		if err == appErrors.ErrNotFound {
			appErrors.RespondError(c, http.StatusNotFound, "Position not found")
			return
		}
		if err == appErrors.ErrPositionAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to update position")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position updated successfully"})
}

func (h *AdminPositionHandler) DeletePosition(c *gin.Context) {
	positionIdParam := c.Param("positionId")
	positionId, err := strconv.Atoi(positionIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid position ID")
		return
	}

	if err := h.positionsService.DeletePosition(c.Request.Context(), uint(positionId)); err != nil {
		if err == appErrors.ErrNotFound {
			appErrors.RespondError(c, http.StatusNotFound, "Position not found")
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to delete position")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Position deleted successfully"})
}
