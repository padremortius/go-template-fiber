package config

import (
	"bytes"
	"errors"
	"fmt"
	"go-template-fiber/internal/common"
	"go-template-fiber/internal/crontab"
	"go-template-fiber/internal/httpserver"
	"go-template-fiber/internal/storage"
	"go-template-fiber/internal/svclogger"
	"net/url"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App     `yaml:"app" json:"app" validate:"required"`
	BaseApp `yaml:"baseApp" json:"baseApp" validate:"required"`
	Crontab crontab.CronOpts   `yaml:"crontab" json:"crontab" validate:"required"`
	HTTP    httpserver.HTTP    `yaml:"http" json:"http" validate:"required"`
	Log     svclogger.Log      `yaml:"logger" json:"logger" validate:"required"`
	Storage storage.StorageCfg `yaml:"storage" json:"storage" validate:"required"`
	Version `json:"version"`
}

var (
	Cfg Config
)

// NewConfig initializes the configuration by reading environment variables
// and a YAML configuration file.
//
// It returns an error if there is an issue reading the environment variables
// or the configuration file.
func NewConfig() error {
	err := cleanenv.ReadEnv(&Cfg)
	if err != nil {
		return err
	}

	appConfigName := fmt.Sprint(Cfg.BaseApp.Name, "-", Cfg.BaseApp.ProfileName, ".yml")

	if Cfg.BaseApp.ProfileName == "dev" {
		if err := cleanenv.ReadConfig(appConfigName, &Cfg); err != nil {
			return err
		}
	} else {
		configURL, _ := url.JoinPath(Cfg.BaseApp.ConfSrvURI, appConfigName)

		data, err := common.GetFileByURL(configURL)
		if err != nil {
			return err
		}

		if err := cleanenv.ParseYAML(bytes.NewBuffer(data), &Cfg); err != nil {
			return err
		}
	}

	// if err = FillPubKey(); err != nil {
	// 	return errors.New("Error read pubkey from oAuth2 server: " + err.Error())
	// }

	if err = ReadPwd(); err != nil {
		return errors.New("Read password error: " + err.Error())
	}

	if err := Cfg.validateConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Config) validateConfig() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}
