package services

import (
	"context"
	"time"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type UserService struct {
	db             *gorm.DB
	userRepository *repositories.UserRepository
	teamRepository *repositories.TeamsRepository
}

func NewUserService(db *gorm.DB, userRepository *repositories.UserRepository, teamRepository *repositories.TeamsRepository) *UserService {
	return &UserService{db: db, userRepository: userRepository, teamRepository: teamRepository}
}

func (s *UserService) GetUserProfile(c context.Context, id uint) (*dtos.UserProfile, error) {
	user, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	userProfile := helpers.MapUserToUserProfile(user)

	return userProfile, nil
}

func (s *UserService) SearchUsers(c context.Context, name *string, teamId *uint, limit, offset int) (*dtos.UserSearchResponse, error) {
	users, totalCount, err := s.userRepository.SearchUsers(s.db.WithContext(c), name, teamId, limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	response := &dtos.UserSearchResponse{
		Users: helpers.MapUsersToUserDataForSearches(users),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}

	return response, nil
}

func (s *UserService) CreateUser(c context.Context, req dtos.CreateOrUpdateUserRequest) error {
	existedUser, err := s.userRepository.FindByEmail(s.db.WithContext(c), req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return appErrors.ErrInternalServerError
	}
	if existedUser != nil {
		return appErrors.ErrEmailAlreadyExists
	}

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

	err = s.db.WithContext(c).Transaction(func(tx *gorm.DB) error {
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

	if err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

func (s *UserService) UpdateUser(c context.Context, id uint, req dtos.CreateOrUpdateUserRequest) error {
	var birthday *time.Time
	if req.Birthday != nil && !req.Birthday.Time.IsZero() {
		birthday = &req.Birthday.Time
	}

	currentUser, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if currentUser.Email != req.Email {
		existedUser, err := s.userRepository.FindByEmail(s.db.WithContext(c), req.Email)
		if err != nil && err != gorm.ErrRecordNotFound {
			return appErrors.ErrInternalServerError
		}
		if existedUser != nil {
			return appErrors.ErrEmailAlreadyExists
		}
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

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := s.userRepository.UpdateUser(tx, user); err != nil {
			return err
		}
		return s.userRepository.UpdateUserSkills(tx, id, userSkills)
	})
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}

func (s *UserService) DeleteUser(c context.Context, id uint) error {
	_, err := s.userRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrUserNotFound
		}
		return appErrors.ErrInternalServerError
	}

	exist, err := s.teamRepository.ExistByLeaderId(s.db.WithContext(c), id)
	if err != nil {
		return appErrors.ErrInternalServerError
	}
	if exist {
		return appErrors.ErrCannotDeleteUserBeingTeamLeader
	}

	if err := s.db.WithContext(c).Delete(&models.User{}, id).Error; err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}
