package services

import (
	"context"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/types"

	"gorm.io/gorm"
)

type UserService struct {
	db             *gorm.DB
	userRepository *repositories.UserRepository
}

func NewUserService(db *gorm.DB, userRepository *repositories.UserRepository) *UserService {
	return &UserService{db: db, userRepository: userRepository}
}

func (s *UserService) GetUserProfile(c context.Context, id uint) (*dtos.UserProfile, error) {
	user, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, err
	}

	var currentTeam dtos.TeamSummary
	if user.CurrentTeam != nil {
		currentTeam = dtos.TeamSummary{
			ID:   user.CurrentTeam.ID,
			Name: user.CurrentTeam.Name,
		}
	}

	projects := make([]dtos.ProjectSummary, 0)
	if len(user.Projects) > 0 {
		for _, project := range user.Projects {
			projects = append(projects, dtos.ProjectSummary{
				ID:           project.ID,
				Name:         project.Name,
				Abbreviation: project.Abbreviation,
				StartDate:    project.StartDate,
				EndDate:      project.EndDate,
			})
		}
	}

	skills := make([]dtos.UserSkillSummary, 0)
	if len(user.Skills) > 0 {
		for _, skill := range user.Skills {
			skills = append(skills, dtos.UserSkillSummary{
				ID:   skill.ID,
				Name: skill.Name,
			})
		}
	}

	userProfile := &dtos.UserProfile{
		ID:          user.ID,
		Name:        user.Name,
		Email:       user.Email,
		Birthday:    &types.Date{Time: *user.Birthday},
		CurrentTeam: &currentTeam,
		Position: dtos.Position{
			ID:           user.Position.ID,
			Name:         user.Position.Name,
			Abbreviation: user.Position.Abbreviation,
		},
		Projects: projects,
		Skills:   skills,
	}

	return userProfile, nil
}

func (s *UserService) SearchUsers(c context.Context, teamId *uint, limit, offset int) (*dtos.UserSearchResponse, error) {
	users, totalCount, err := s.userRepository.SearchUsers(s.db.WithContext(c), teamId, limit, offset)
	if err != nil {
		return nil, err
	}

	userDtos := make([]dtos.UserDataForSearch, 0, len(users))
	if len(users) > 0 {
		for _, user := range users {
			userDtos = append(
				userDtos,
				dtos.UserDataForSearch{
					ID:    user.ID,
					Name:  user.Name,
					Email: user.Email,
				})
		}
	}

	response := &dtos.UserSearchResponse{
		Users: userDtos,
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}
