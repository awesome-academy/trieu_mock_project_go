package services

import (
	"context"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type TeamsService struct {
	db                   *gorm.DB
	teamRepository       *repositories.TeamsRepository
	teamMemberRepository *repositories.TeamMemberRepository
}

func NewTeamsService(db *gorm.DB, teamRepository *repositories.TeamsRepository, teamMemberRepository *repositories.TeamMemberRepository) *TeamsService {
	return &TeamsService{db: db, teamRepository: teamRepository, teamMemberRepository: teamMemberRepository}
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

func (s *TeamsService) GetTeamDetails(c context.Context, id uint) (*dtos.Team, error) {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, err
	}

	return helpers.MapTeamToTeamDto(team), nil
}

func (s *TeamsService) GetTeamMembers(c context.Context, teamID uint, limit, offset int) (*dtos.ListTeamMembersResponse, error) {
	members, err := s.teamMemberRepository.FindMembersByTeamID(s.db.WithContext(c), teamID, limit, offset)
	if err != nil {
		return nil, err
	}

	totalCount, err := s.teamMemberRepository.CountMembersByTeamID(s.db.WithContext(c), teamID)
	if err != nil {
		return nil, err
	}

	response := &dtos.ListTeamMembersResponse{
		Members: helpers.MapTeamMembersToTeamMemberSummaries(members),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *TeamsService) GetAllTeamsSummary(c context.Context) []dtos.TeamSummary {
	teams, err := s.teamRepository.FindAllTeamsSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.TeamSummary{}
	}

	return helpers.MapTeamsToTeamSummaries(teams)
}

func (r *TeamsService) extractTeamMembersFromTeam(team models.Team) []dtos.UserSummary {
	teamMembers := make([]dtos.UserSummary, 0)
	if len(team.Members) > 0 {
		for _, member := range team.Members {
			teamMembers = append(teamMembers, dtos.UserSummary{
				ID:   member.ID,
				Name: member.Name,
			})
		}
	}
	return teamMembers
}

func (r *TeamsService) extractProjectsFromTeam(team models.Team) []dtos.ProjectSummary {
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
	return projects
}
