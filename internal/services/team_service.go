package services

import (
	"context"
	"trieu_mock_project_go/helpers"
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

	teamDtos := helpers.MapTeamsToTeamDtos(teams)

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
