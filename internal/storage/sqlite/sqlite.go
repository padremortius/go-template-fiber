package sqlite

import (
	"context"
	"os"
	"path/filepath"

	"github.com/padremortius/go-template-fiber/internal/svclogger"

	"github.com/glebarez/sqlite"

	"gorm.io/gorm"
)

type (
	Storage struct {
		db  *gorm.DB
		log svclogger.Log
	}
)

func New(aCtx context.Context, path string, log *svclogger.Log) (*Storage, error) {
	dbPath := filepath.Dir(path)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if err := os.Mkdir(dbPath, os.ModePerm); err != nil {
			return nil, err
		}
	}
	log.Logger.Debug().Msgf("Start init new storage at path: %s", path)
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &Storage{db: db, log: *log}, nil
}
