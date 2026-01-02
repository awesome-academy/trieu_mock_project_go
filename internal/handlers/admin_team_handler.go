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

type AdminTeamHandler struct {
	teamService *services.TeamsService
	userService *services.UserService
}

func NewAdminTeamHandler(teamService *services.TeamsService, userService *services.UserService) *AdminTeamHandler {
	return &AdminTeamHandler{teamService: teamService, userService: userService}
}

func (h *AdminTeamHandler) ListTeamPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_teams.html", gin.H{
		"title":     "Admin Teams Management",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminTeamHandler) TeamSearchPartial(c *gin.Context) {
	templateName := "partials/admin_teams_search.html"
	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid query parameters")
		return
	}

	resp, err := h.teamService.ListTeams(c.Request.Context(), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load teams")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"teams": resp.Teams,
		"page":  resp.Page,
	})
}

func (h *AdminTeamHandler) CreateTeamPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_team_create.html", gin.H{
		"title":     "Create Team",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminTeamHandler) CreateTeam(c *gin.Context) {
	var request dtos.CreateOrUpdateTeamRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.teamService.CreateTeam(c.Request.Context(), request); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to create team")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team created successfully"})
}

func (h *AdminTeamHandler) EditTeamPage(c *gin.Context) {
	templateName := "pages/admin_team_edit.html"
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid team ID")
		return
	}

	team, err := h.teamService.GetTeamDetails(c.Request.Context(), uint(teamId))
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Team not found")
		return
	}

	historyResp, err := h.teamService.GetTeamMemberHistory(c.Request.Context(), uint(teamId), 10, 0)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load team member history")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"title":     "Edit Team",
		"team":      team,
		"history":   historyResp.History,
		"page":      historyResp.Page,
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminTeamHandler) TeamMemberHistoryPartial(c *gin.Context) {
	templateName := "partials/admin_team_member_history.html"
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid team ID")
		return
	}

	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid query parameters")
		return
	}

	if requestQuery.Limit == 0 {
		requestQuery.Limit = 10
	}

	resp, err := h.teamService.GetTeamMemberHistory(c.Request.Context(), uint(teamId), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load team member history")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"history": resp.History,
		"page":    resp.Page,
	})
}

func (h *AdminTeamHandler) UpdateTeam(c *gin.Context) {
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var request dtos.CreateOrUpdateTeamRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.teamService.UpdateTeam(c.Request.Context(), uint(teamId), request); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to update team")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team updated successfully"})
}

func (h *AdminTeamHandler) DeleteTeam(c *gin.Context) {
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	if err := h.teamService.DeleteTeam(c.Request.Context(), uint(teamId)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to delete team")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

func (h *AdminTeamHandler) AddMember(c *gin.Context) {
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	var request dtos.AddMemberRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.teamService.AddMemberToTeam(c.Request.Context(), uint(teamId), request.UserID); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to add member")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member added successfully"})
}

func (h *AdminTeamHandler) RemoveMember(c *gin.Context) {
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid team ID")
		return
	}

	userIdParam := c.Param("userId")
	userId, err := strconv.Atoi(userIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := h.teamService.RemoveMemberFromTeam(c.Request.Context(), uint(teamId), uint(userId)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to remove member")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}
