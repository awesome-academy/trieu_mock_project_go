package services

import (
	"context"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"

	"gorm.io/gorm"
)

type TeamsService struct {
	db             *gorm.DB
	teamRepository *repositories.TeamsRepository
}

func NewTeamsService(db *gorm.DB, teamRepository *repositories.TeamsRepository) *TeamsService {
	return &TeamsService{db: db, teamRepository: teamRepository}
}

func (s *TeamsService) ListTeams(c context.Context, limit, offset int) (*dtos.ListTeamsResponse, error) {
	teams, err := s.teamRepository.ListTeams(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.teamRepository.CountTeams(s.db.WithContext(c))
	if err != nil {
		return nil, err
	}

	teamDtos := make([]dtos.Team, 0, len(teams))
	for _, team := range teams {

		teamMembers := make([]dtos.UserSummary, 0)
		if len(team.Members) > 0 {
			for _, member := range team.Members {
				teamMembers = append(teamMembers, dtos.UserSummary{
					ID:   member.ID,
					Name: member.Name,
				})
			}
		}

		projects := make([]dtos.ProjectSummary, 0, len(team.Projects))
		if len(team.Projects) > 0 {
				for _, project := range team.Projects {
					projects = append(projects, dtos.ProjectSummary{
						ID:           project.ID,
						Name:         project.Name,
						Abbreviation: project.Abbreviation,
						StartDate:    project.StartDate,
						EndDate:      project.EndDate,
					})
				}
		}

		teamDtos = append(teamDtos, dtos.Team{
			ID:          team.ID,
			Name:        team.Name,
			Description: team.Description,
			CreatedAt:   team.CreatedAt,
			UpdatedAt:   team.UpdatedAt,

			Leader: dtos.UserSummary{
				ID:   team.LeaderID,
				Name: team.Leader.Name,
			},
			Members:  teamMembers,
			Projects: projects,
		})
	}

	response := &dtos.ListTeamsResponse{
		Teams: teamDtos,
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}
