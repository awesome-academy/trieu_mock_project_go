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

type AdminProjectHandler struct {
	projectService *services.ProjectService
	teamService    *services.TeamsService
	userService    *services.UserService
}

func NewAdminProjectHandler(
	projectService *services.ProjectService,
	teamService *services.TeamsService,
	userService *services.UserService,
) *AdminProjectHandler {
	return &AdminProjectHandler{
		projectService: projectService,
		teamService:    teamService,
		userService:    userService,
	}
}

func (h *AdminProjectHandler) ListProjectPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_projects.html", gin.H{
		"title":     "Admin Projects Management",
		"teams":     c.MustGet("teams"),
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminProjectHandler) ProjectSearchPartial(c *gin.Context) {
	templateName := "partials/admin_projects_search.html"
	var requestQuery dtos.ProjectSearchRequest
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid query parameters")
		return
	}

	resp, err := h.projectService.SearchProjects(c.Request.Context(), requestQuery.TeamID, requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load projects")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"projects": resp.Projects,
		"page":     resp.Page,
	})
}

func (h *AdminProjectHandler) CreateProjectPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_project_create.html", gin.H{
		"title":     "Create Project",
		"teams":     c.MustGet("teams"),
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminProjectHandler) CreateProject(c *gin.Context) {
	var request dtos.CreateOrUpdateProjectRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.projectService.CreateProject(c.Request.Context(), request); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to create project")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project created successfully"})
}

func (h *AdminProjectHandler) EditProjectPage(c *gin.Context) {
	templateName := "pages/admin_project_edit.html"
	projectIdParam := c.Param("projectId")
	projectId, err := strconv.Atoi(projectIdParam)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid project ID")
		return
	}

	project, err := h.projectService.GetProjectByID(c.Request.Context(), uint(projectId))
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load project")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"title":     "Edit Project",
		"project":   project,
		"teams":     c.MustGet("teams"),
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminProjectHandler) UpdateProject(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	projectId, err := strconv.Atoi(projectIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	var request dtos.CreateOrUpdateProjectRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.projectService.UpdateProject(c.Request.Context(), uint(projectId), request); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to update project")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project updated successfully"})
}

func (h *AdminProjectHandler) DeleteProject(c *gin.Context) {
	projectIdParam := c.Param("projectId")
	projectId, err := strconv.Atoi(projectIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid project ID")
		return
	}

	if err := h.projectService.DeleteProject(c.Request.Context(), uint(projectId)); err != nil {
		appErrors.RespondCustomError(c, err, "Failed to delete project")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
