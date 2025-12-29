package services

import (
	"context"
	"time"
	"trieu_mock_project_go/internal/dtos"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"
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

	var currentTeam dtos.PositionSummary
	if user.CurrentTeam != nil {
		currentTeam = dtos.PositionSummary{
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

	userSkills := make([]dtos.UserSkillSummary, 0)
	if len(user.UserSkill) > 0 {
		for _, skill := range user.UserSkill {
			userSkills = append(userSkills, dtos.UserSkillSummary{
				ID:             skill.Skill.ID,
				Name:           skill.Skill.Name,
				Level:          skill.Level,
				UsedYearNumber: skill.UsedYearNumber,
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
		Skills:   userSkills,
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

func (s *UserService) CreateUser(c context.Context, req dtos.CreateOrUpdateUserRequest) error {
	var birthday *time.Time
	if req.Birthday != nil && !req.Birthday.Time.IsZero() {
		birthday = &req.Birthday.Time
	}

	user := &models.User{
		Name:          req.Name,
		Email:         req.Email,
		Birthday:      birthday,
		PositionID:    req.PositionID,
		CurrentTeamID: req.TeamID,
	}

	s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
		err := s.userRepository.CreateUser(tx, user)
		if err != nil {
			return err
		}

		userSkills := make([]models.UserSkill, 0, len(req.Skills))
		for _, sReq := range req.Skills {
			userSkills = append(userSkills, models.UserSkill{
				UserID:         user.ID,
				SkillID:        sReq.ID,
				Level:          sReq.Level,
				UsedYearNumber: sReq.UsedYearNumber,
			})
		}

		return s.userRepository.CreateUserSkills(tx, userSkills)
	})

	return nil
}

func (s *UserService) UpdateUser(c context.Context, id uint, req dtos.CreateOrUpdateUserRequest) error {
	var birthday *time.Time
	if req.Birthday != nil && !req.Birthday.Time.IsZero() {
		birthday = &req.Birthday.Time
	}

	user := &models.User{
		ID:            id,
		Name:          req.Name,
		Email:         req.Email,
		Birthday:      birthday,
		PositionID:    req.PositionID,
		CurrentTeamID: req.TeamID,
	}

	userSkills := make([]models.UserSkill, 0, len(req.Skills))
	for _, sReq := range req.Skills {
		userSkills = append(userSkills, models.UserSkill{
			UserID:         id,
			SkillID:        sReq.ID,
			Level:          sReq.Level,
			UsedYearNumber: sReq.UsedYearNumber,
		})
	}

	return s.userRepository.UpdateUser(s.db.WithContext(c), user, userSkills)
}

func (s *UserService) DeleteUser(c context.Context, id uint) error {
	if err := s.db.WithContext(c).Delete(&models.User{}, id).Error; err != nil {
		return err
	}

	return nil
}
