package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type ActivityLogRepository struct {
}

func NewActivityLogRepository() *ActivityLogRepository {
	return &ActivityLogRepository{}
}

func (r *ActivityLogRepository) Create(db *gorm.DB, log *models.ActivityLog) error {
	return db.Create(log).Error
}

func (r *ActivityLogRepository) SearchActivityLogs(db *gorm.DB, limit, offset int) ([]models.ActivityLog, int64, error) {
	var logs []models.ActivityLog
	query := db.Model(&models.ActivityLog{}).Preload("User")

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, count, nil
}

func (r *ActivityLogRepository) FindByID(db *gorm.DB, id uint) (*models.ActivityLog, error) {
	var log models.ActivityLog
	result := db.First(&log, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &log, nil
}

func (r *ActivityLogRepository) Delete(db *gorm.DB, id uint) error {
	return db.Delete(&models.ActivityLog{}, id).Error
}

func (r *ActivityLogRepository) CreateInBatches(db *gorm.DB, logs []models.ActivityLog, batchSize int) error {
	return db.CreateInBatches(logs, batchSize).Error
}
