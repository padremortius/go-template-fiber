package storage

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/padremortius/go-template-fiber/pkgs/svclogger"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var ErrInvalidUserType = errors.New("invalid user type")

type (
	StorageCfg struct {
		DBType string `yaml:"dbType" json:"dbType" validate:"required"`
		Path   string `yaml:"path" json:"path" validate:"required"`
	}
)

type StorageInterface interface {
	InitDB(ctx context.Context) error
	AddUser(ctx context.Context, user interface{}) error
	DeleteUser(ctx context.Context, tgUserID int64) error
	GetActiveUserByCurrentDay(ctx context.Context) (int64, error)
	GetListActiveUsers(ctx context.Context) (interface{}, error)
	Close() error
}

type (
	Storage struct {
		db  *gorm.DB
		log *svclogger.Log
	}
)

func New(aCtx context.Context, path string, log *svclogger.Log) (StorageInterface, error) {
	dbPath := filepath.Dir(path)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.Mkdir(dbPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	log.Debugf("Start init new storage at path: %v", path)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: db, log: log}, nil
}

func (s *Storage) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
