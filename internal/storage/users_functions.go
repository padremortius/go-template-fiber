package storage

import (
	"context"
	"time"

	"github.com/padremortius/go-template-fiber/internal/structs/models"

	"gorm.io/gorm/clause"
)

func (s *Storage) InitDB(ctx context.Context) error {
	if err := s.db.WithContext(ctx).AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddUser(ctx context.Context, user interface{}) error {
	u, ok := user.(models.User)
	if !ok {
		return ErrInvalidUserType
	}
	return s.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"active", "create_date_time"}),
		}).Create(u).Error
}

func (s *Storage) DeleteUser(ctx context.Context, tgUserID int64) error {
	s.log.Debugf("Delete user: %v", tgUserID)
	return s.db.WithContext(ctx).Model(
		&models.User{},
	).Where(
		"user_id = ?", tgUserID,
	).Updates(
		models.User{Active: 0, UpdateDateTime: time.Now()},
	).Error
}

func (s *Storage) GetActiveUserByCurrentDay(ctx context.Context) (int64, error) {
	var res int64
	err := s.db.WithContext(ctx).Model(
		&models.User{},
	).Where(
		"active = ? and update_date_time >= ?",
		1,
		time.Now().Format("2006-01-02")+" 00:00:00",
		1).Count(&res).Error
	return res, err
}

func (s *Storage) GetListActiveUsers(ctx context.Context) (interface{}, error) {
	var res []models.User
	err := s.db.WithContext(ctx).Model(
		&models.User{},
	).Where(
		"active = ?",
		1,
	).Find(&res).Error
	return res, err
}
