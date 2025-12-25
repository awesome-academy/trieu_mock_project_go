package handlers

import (
	"net/http"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type TeamsHandler struct {
	teamsService *services.TeamsService
}

func NewTeamsHandler(teamsService *services.TeamsService) *TeamsHandler {
	return &TeamsHandler{teamsService: teamsService}
}

func (h *TeamsHandler) TeamsPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/teams.html", gin.H{
		"title": "Teams",
	})
}

func (h *TeamsHandler) ListTeams(c *gin.Context) {
	var query dtos.ListTeamsRequestQuery
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&query)) {
		return
	}

	resp, err := h.teamsService.ListTeams(c.Request.Context(), query.Limit, query.Offset)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to list teams")
		return
	}

	c.JSON(http.StatusOK, resp)
}
