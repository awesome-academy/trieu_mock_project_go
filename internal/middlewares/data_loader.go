package middlewares

import (
	"trieu_mock_project_go/internal/services"

	"github.com/gin-gonic/gin"
)

type DataLoader struct {
	teamService     *services.TeamsService
	positionService *services.PositionService
	skillService    *services.SkillService
}

func NewDataLoader(
	teamService *services.TeamsService,
	positionService *services.PositionService,
	skillService *services.SkillService,
) *DataLoader {
	return &DataLoader{
		teamService:     teamService,
		positionService: positionService,
		skillService:    skillService,
	}
}

func (dl *DataLoader) LoadTeams() gin.HandlerFunc {
	return func(c *gin.Context) {
		teams := dl.teamService.GetAllTeamsSummary(c.Request.Context())
		c.Set("teams", teams)
		c.Next()
	}
}

func (dl *DataLoader) LoadTeamPositionSkill() gin.HandlerFunc {
	return func(c *gin.Context) {
		teams := dl.teamService.GetAllTeamsSummary(c.Request.Context())
		positions := dl.positionService.GetAllPositionsSummary(c.Request.Context())
		skills := dl.skillService.GetAllSkillsSummary(c.Request.Context())
		c.Set("teams", teams)
		c.Set("positions", positions)
		c.Set("skills", skills)
		c.Next()
	}
}
