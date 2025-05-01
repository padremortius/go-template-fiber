package storage

import (
	"time"

	"github.com/padremortius/go-template-fiber/internal/structs/models"

	"gorm.io/gorm/clause"
)

func (s *Storage) InitDB() error {
	if err := s.db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	return nil
}

func (s *Storage) AddUser(u models.User) error {
	return s.db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"active", "create_date_time"}),
		}).Create(u).Error
}

func (s *Storage) DeleteUser(tgUserId int64) error {
	s.log.Logger.Debug().Msgf("Delete user: %v", tgUserId)
	return s.db.Model(
		&models.User{},
	).Where(
		"user_id = ?", tgUserId,
	).Updates(
		models.User{Active: 0, UpdateDateTime: time.Now()},
	).Error
}

func (s *Storage) GetActiveUserByCurrentDay() (int64, error) {
	var res int64
	err := s.db.Model(
		&models.User{},
	).Where(
		"active = ? and update_date_time >= ?",
		1,
		time.Now().Format("2006-01-02")+" 00:00:00",
		1).Count(&res).Error
	return res, err
}

func (s *Storage) GetListActiveUsers() ([]models.User, error) {
	var res []models.User
	err := s.db.Model(
		&models.User{},
	).Where(
		"active = ?",
		1,
	).Find(&res).Error
	return res, err
}
