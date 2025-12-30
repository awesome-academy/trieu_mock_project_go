package handlers

import (
	"net/http"
	"strconv"
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

func (h *TeamsHandler) TeamDetailsPageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/team_details.html", gin.H{
		"title": "Team Details",
	})
}

func (h *TeamsHandler) ListTeams(c *gin.Context) {
	var query dtos.PaginationRequestQuery
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&query)) {
		return
	}

	teamsResp, err := h.teamsService.ListTeams(c.Request.Context(), query.Limit, query.Offset)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to list teams")
		return
	}

	c.JSON(http.StatusOK, teamsResp)
}

func (h *TeamsHandler) GetTeamDetails(c *gin.Context) {
	teamIdParam := c.Param("id")

	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	resp, err := h.teamsService.GetTeamDetails(c.Request.Context(), uint(teamId))
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to get team details")
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *TeamsHandler) GetTeamMembers(c *gin.Context) {
	teamIdParam := c.Param("id")

	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var query dtos.PaginationRequestQuery
	if appErrors.HandleBindError(c, c.ShouldBindQuery(&query)) {
		return
	}

	resp, err := h.teamsService.GetTeamMembers(c.Request.Context(), uint(teamId), query.Limit, query.Offset)
	if err != nil {
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to get team members")
		return
	}

	c.JSON(http.StatusOK, resp)
}
