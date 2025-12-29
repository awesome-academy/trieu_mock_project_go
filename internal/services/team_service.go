package services

import (
	"context"
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

	teamDtos := make([]dtos.Team, 0, len(teams))
	for _, team := range teams {

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
			Members:  s.extractTeamMembersFromTeam(team),
			Projects: s.extractProjectsFromTeam(team),
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

func (s *TeamsService) GetTeamDetails(c context.Context, id uint) (*dtos.Team, error) {
	team, err := s.teamRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, err
	}

	teamDto := &dtos.Team{
		ID:          team.ID,
		Name:        team.Name,
		Description: team.Description,
		CreatedAt:   team.CreatedAt,
		UpdatedAt:   team.UpdatedAt,

		Leader: dtos.UserSummary{
			ID:   team.LeaderID,
			Name: team.Leader.Name,
		},
		Members:  s.extractTeamMembersFromTeam(*team),
		Projects: s.extractProjectsFromTeam(*team),
	}

	return teamDto, nil
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

	memberDtos := make([]dtos.TeamMemberSummary, 0, len(members))
	if len(members) > 0 {
		for _, member := range members {
			memberDtos = append(memberDtos, dtos.TeamMemberSummary{
				ID:       member.User.ID,
				Name:     member.User.Name,
				Email:    member.User.Email,
				JoinedAt: member.JoinedAt,
			})
		}
	}

	response := &dtos.ListTeamMembersResponse{
		Members: memberDtos,
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *TeamsService) GetAllTeamsSummary(c context.Context) []dtos.PositionSummary {
	teams, err := s.teamRepository.FindAllTeamsSummary(s.db.WithContext(c))
	if err != nil {
		return []dtos.PositionSummary{}
	}

	teamDtos := make([]dtos.PositionSummary, 0, len(teams))
	for _, team := range teams {
		teamDtos = append(teamDtos, dtos.PositionSummary{
			ID:   team.ID,
			Name: team.Name,
		})
	}
	return teamDtos
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
