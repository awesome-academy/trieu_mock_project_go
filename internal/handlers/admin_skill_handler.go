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

type AdminSkillHandler struct {
	skillService *services.SkillService
}

func NewAdminSkillHandler(skillService *services.SkillService) *AdminSkillHandler {
	return &AdminSkillHandler{skillService: skillService}
}

func (h *AdminSkillHandler) ListSkillPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_skills.html", gin.H{
		"title":     "Admin Skills Management",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminSkillHandler) SkillSearchPartial(c *gin.Context) {
	templateName := "partials/admin_skills_search.html"
	var requestQuery dtos.PaginationRequestQuery
	if err := c.ShouldBindQuery(&requestQuery); err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid query parameters")
		return
	}

	resp, err := h.skillService.SearchSkills(c.Request.Context(), requestQuery.Limit, requestQuery.Offset)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load skills")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"skills": resp.Skills,
		"page":   resp.Page,
	})
}

func (h *AdminSkillHandler) CreateSkillPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_skill_create.html", gin.H{
		"title":     "Create Skill",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminSkillHandler) CreateSkill(c *gin.Context) {
	var request dtos.CreateOrUpdateSkillRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.skillService.CreateSkill(c.Request.Context(), request); err != nil {
		if err == appErrors.ErrSkillAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to create skill")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill created successfully"})
}

func (h *AdminSkillHandler) EditSkillPage(c *gin.Context) {
	templateName := "pages/admin_skill_edit.html"
	skillIdParam := c.Param("skillId")
	skillId, err := strconv.Atoi(skillIdParam)
	if err != nil {
		appErrors.RespondPageError(c, http.StatusBadRequest, templateName, "Invalid skill ID")
		return
	}

	skill, err := h.skillService.GetSkillByID(c.Request.Context(), uint(skillId))
	if err != nil {
		appErrors.RespondPageError(c, http.StatusInternalServerError, templateName, "Failed to load skill")
		return
	}

	c.HTML(http.StatusOK, templateName, gin.H{
		"title":     "Edit Skill",
		"skill":     skill,
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminSkillHandler) UpdateSkill(c *gin.Context) {
	skillIdParam := c.Param("skillId")
	skillId, err := strconv.Atoi(skillIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid skill ID")
		return
	}

	var request dtos.CreateOrUpdateSkillRequest
	if appErrors.HandleBindError(c, c.ShouldBindJSON(&request)) {
		return
	}

	if err := h.skillService.UpdateSkill(c.Request.Context(), uint(skillId), request); err != nil {
		if err == appErrors.ErrNotFound {
			appErrors.RespondError(c, http.StatusNotFound, "Skill not found")
			return
		}
		if err == appErrors.ErrSkillAlreadyExists {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to update skill")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill updated successfully"})
}

func (h *AdminSkillHandler) DeleteSkill(c *gin.Context) {
	skillIdParam := c.Param("skillId")
	skillId, err := strconv.Atoi(skillIdParam)
	if err != nil {
		appErrors.RespondError(c, http.StatusBadRequest, "Invalid skill ID")
		return
	}

	if err := h.skillService.DeleteSkill(c.Request.Context(), uint(skillId)); err != nil {
		if err == appErrors.ErrNotFound {
			appErrors.RespondError(c, http.StatusNotFound, "Skill not found")
			return
		}
		if err == appErrors.ErrSkillInUse {
			appErrors.RespondError(c, http.StatusConflict, err.Error())
			return
		}
		appErrors.RespondError(c, http.StatusInternalServerError, "Failed to delete skill")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}
