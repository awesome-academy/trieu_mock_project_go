package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"time"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type AdminExportCsvHandler struct {
	userService     *services.UserService
	positionService *services.PositionService
	projectService  *services.ProjectService
	skillService    *services.SkillService
	teamService     *services.TeamsService
}

func NewAdminExportCsvHandler(userService *services.UserService, positionService *services.PositionService, projectService *services.ProjectService, skillService *services.SkillService, teamService *services.TeamsService) *AdminExportCsvHandler {
	return &AdminExportCsvHandler{
		userService:     userService,
		positionService: positionService,
		projectService:  projectService,
		skillService:    skillService,
		teamService:     teamService,
	}
}

func (h *AdminExportCsvHandler) ExportCSVPage(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/admin_export_csv.html", gin.H{
		"title":     "Export Data as CSV",
		"csrfToken": csrf.GetToken(c),
	})
}

func (h *AdminExportCsvHandler) ExportCSV(c *gin.Context) {
	templateName := "pages/admin_export_csv.html"
	var req dtos.ExportRequest
	if err := c.ShouldBind(&req); err != nil {
		appErrors.RespondPageErrorWithCSRF(c, http.StatusBadRequest, templateName, "Invalid request data")
		return
	}

	var csvData [][]string
	var err error

	fileName := req.Type + "_" + time.Now().Format("20060102_150405") + ".csv"
	switch req.Type {
	case "user":
		csvData, err = h.userService.ExportUsersToCSV(c)
	case "position":
		csvData, err = h.positionService.ExportPositionsToCSV(c)
	case "project":
		csvData, err = h.projectService.ExportProjectsToCSV(c)
	case "skill":
		csvData, err = h.skillService.ExportSkillsToCSV(c)
	case "team":
		csvData, err = h.teamService.ExportTeamsToCSV(c)
	}

	if err != nil {
		appErrors.RespondPageErrorWithCSRF(c, http.StatusInternalServerError, templateName, "Failed to fetch data")
		return
	}

	b := &bytes.Buffer{}
	w := csv.NewWriter(b)
	if err := w.WriteAll(csvData); err != nil {
		appErrors.RespondPageErrorWithCSRF(c, http.StatusInternalServerError, templateName, "Failed to generate CSV")
		return
	}

	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Data(http.StatusOK, "text/csv", b.Bytes())
}
