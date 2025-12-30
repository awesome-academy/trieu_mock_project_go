package handlers

import (
	"net/http"
	"strconv"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"gorm.io/gorm"
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
	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		c.HTML(http.StatusBadRequest, "partials/admin_teams_search.html", gin.H{
			"error": "Invalid query parameters",
		})
		return
	}

	resp, err := h.teamService.ListTeams(c.Request.Context(), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "partials/admin_teams_search.html", gin.H{
			"error": "Failed to load teams",
		})
		return
	}

	c.HTML(http.StatusOK, "partials/admin_teams_search.html", gin.H{
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
		if err == gorm.ErrRecordNotFound {
			appErrors.RespondError(c, http.StatusBadRequest, "Leader user not found")
			return
		}
		if err == appErrors.ErrTeamLeaderAlreadyInAnotherTeam {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		if err == appErrors.ErrTeamAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to create team")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Team created successfully"})
}

func (h *AdminTeamHandler) EditTeamPage(c *gin.Context) {
	teamIdParam := c.Param("teamId")
	teamId, err := strconv.Atoi(teamIdParam)
	if err != nil {
		c.HTML(http.StatusBadRequest, "pages/admin_team_edit.html", gin.H{
			"title": "Edit Team",
			"error": "Invalid team ID",
		})
		return
	}

	team, err := h.teamService.GetTeamDetails(c.Request.Context(), uint(teamId))
	if err != nil {
		c.HTML(http.StatusInternalServerError, "pages/admin_team_edit.html", gin.H{
			"title": "Edit Team",
			"error": "Team not found",
		})
		return
	}

	c.HTML(http.StatusOK, "pages/admin_team_edit.html", gin.H{
		"title":     "Edit Team",
		"team":      team,
		"csrfToken": csrf.GetToken(c),
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
		if err == appErrors.ErrTeamAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to update team")
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
		if err == appErrors.ErrNotFound {
			appErrors.RespondError(c, http.StatusNotFound, "Team not found")
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to delete team")
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
		if err == appErrors.ErrTeamNotFound ||
			err == appErrors.ErrUserNotFound ||
			err == appErrors.ErrUserAlreadyInTeam {
			appErrors.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to add member")
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
		if err == appErrors.ErrTeamNotFound ||
			err == appErrors.ErrUserNotFound ||
			err == appErrors.ErrUserNotInTeam {
			appErrors.RespondError(c, http.StatusBadRequest, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to remove member")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Member removed successfully"})
}
