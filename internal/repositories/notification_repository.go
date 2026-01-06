package repositories

import (
	"trieu_mock_project_go/models"

	"gorm.io/gorm"
)

type NotificationRepository struct {
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(db *gorm.DB, notification *models.Notification) error {
	return db.Create(notification).Error
}

func (r *NotificationRepository) CreateInBatches(db *gorm.DB, notifications []models.Notification, batchSize int) error {
	return db.CreateInBatches(notifications, batchSize).Error
}

func (r *NotificationRepository) FindByUserID(db *gorm.DB, userID uint, limit, offset int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := db.Model(&models.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Limit(limit).Offset(offset).Find(&notifications).Error; err != nil {
		return nil, 0, err
	}

	return notifications, total, nil
}

func (r *NotificationRepository) CountUnreadByUserID(db *gorm.DB, userID uint) (int64, error) {
	var count int64
	if err := db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *NotificationRepository) UpdateNotificationAsRead(db *gorm.DB, userID uint, notificationID uint) error {
	return db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) UpdateAllNotificationsAsReadByUserID(db *gorm.DB, userID uint) error {
	return db.Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (r *NotificationRepository) DeleteByUserIDAndNotificationID(db *gorm.DB, userID uint, notificationID uint) error {
	return db.Where("id = ? AND user_id = ?", notificationID, userID).
		Delete(&models.Notification{}).Error
}
