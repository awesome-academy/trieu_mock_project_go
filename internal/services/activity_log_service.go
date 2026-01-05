package services

import (
	"context"
	"fmt"
	"log"
	"trieu_mock_project_go/helpers"
	"trieu_mock_project_go/internal/dtos"
	appErrors "trieu_mock_project_go/internal/errors"
	"trieu_mock_project_go/internal/repositories"
	"trieu_mock_project_go/internal/types"
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type ActivityLogService struct {
	db                    *gorm.DB
	activityLogRepository *repositories.ActivityLogRepository
}

func NewActivityLogService(db *gorm.DB, activityLogRepository *repositories.ActivityLogRepository) *ActivityLogService {
	return &ActivityLogService{db: db, activityLogRepository: activityLogRepository}
}

func (s *ActivityLogService) LogActivity(c context.Context, action types.ActivityLog, params ...interface{}) *appErrors.AppError {
	return s.LogActivityDb(c, s.db.WithContext(c), action, params...)
}

func (s *ActivityLogService) createLogActivityModel(c context.Context, action types.ActivityLog, params ...interface{}) (*models.ActivityLog, *appErrors.AppError) {
	var userId uint
	var userEmail string
	var allParams []interface{}

	isAuthAction := action.Value == types.UserSignIn.Value ||
		action.Value == types.AdminSignIn.Value ||
		action.Value == types.UserSignOut.Value ||
		action.Value == types.AdminSignOut.Value

	if isAuthAction {
		if len(params) >= 2 {
			if id, ok := params[0].(uint); ok {
				userId = id
			}
			if email, ok := params[1].(string); ok {
				userEmail = email
			}
		}
		allParams = params
	} else {
		var ok bool
		userId, ok = c.Value("user_id").(uint)
		log.Default().Println("User ID from context:", userId)
		if !ok {
			log.Default().Println("Failed to get user_id from context")
			return nil, appErrors.ErrInternalServerError
		}
		userEmail, ok = c.Value("email").(string)
		log.Default().Println("User Email from context:", userEmail)
		if !ok {
			log.Default().Println("Failed to get email from context")
			return nil, appErrors.ErrInternalServerError
		}
		allParams = append([]interface{}{userId, userEmail}, params...)
	}

	var description *string
	if action.DescriptionFormat != "" {
		desc := fmt.Sprintf(action.DescriptionFormat, allParams...)
		description = &desc
	}

	activityLog := &models.ActivityLog{
		UserID:      userId,
		Action:      action.Value,
		Description: description,
	}

	return activityLog, nil
}

func (s *ActivityLogService) LogActivityDb(c context.Context, db *gorm.DB, action types.ActivityLog, params ...interface{}) *appErrors.AppError {
	activityLog, err := s.createLogActivityModel(c, action, params...)
	if err != nil {
		return err
	}

	if er := s.activityLogRepository.Create(db, activityLog); er != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

func (s *ActivityLogService) SearchActivityLogs(c context.Context, limit, offset int) (*dtos.ActivityLogSearchResponse, *appErrors.AppError) {
	logs, totalCount, err := s.activityLogRepository.SearchActivityLogs(s.db.WithContext(c), limit, offset)
	if err != nil {
		return nil, appErrors.ErrInternalServerError
	}

	return &dtos.ActivityLogSearchResponse{
		Logs: helpers.MapActivityLogsToActivityLogSummaries(logs),
		Page: dtos.PaginationResponse{
			Limit:  limit,
			Offset: offset,
			Total:  totalCount,
		},
	}, nil
}

func (s *ActivityLogService) DeleteActivityLog(c context.Context, id uint) *appErrors.AppError {
	_, err := s.activityLogRepository.FindByID(s.db.WithContext(c), id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return appErrors.ErrActivityLogNotFound
		}
		return appErrors.ErrInternalServerError
	}

	if err := s.activityLogRepository.Delete(s.db.WithContext(c), id); err != nil {
		return appErrors.ErrInternalServerError
	}

	return nil
}

func (s *ActivityLogService) createInBatches(db *gorm.DB, activityLogs []models.ActivityLog, batchSize int) *appErrors.AppError {
	if err := s.activityLogRepository.CreateInBatches(db, activityLogs, batchSize); err != nil {
		return appErrors.ErrInternalServerError
	}
	return nil
}
