package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/padremortius/go-template-fiber/internal/crontab"
	"github.com/padremortius/go-template-fiber/internal/storage"
	"github.com/padremortius/go-template-fiber/pkgs/common"
	"github.com/padremortius/go-template-fiber/pkgs/httpserver"
	"github.com/padremortius/go-template-fiber/pkgs/svclogger"
)

type (
	App struct {
		AuthSrvPassword string `yaml:"-" json:"-" validate:"required"`
	}

	pwdData map[string]string

	Config struct {
		App     `yaml:"app" json:"app" validate:"required"`
		BaseApp `yaml:"baseApp" json:"baseApp" validate:"required"`
		Crontab crontab.CronOpts   `yaml:"crontab" json:"crontab" validate:"required"`
		HTTP    httpserver.HTTP    `yaml:"http" json:"http" validate:"required"`
		Log     svclogger.Log      `yaml:"logger" json:"logger" validate:"required"`
		Storage storage.StorageCfg `yaml:"storage" json:"storage" validate:"required"`
		Version `json:"version"`
	}
)

func (c *Config) ReadBaseConfig() error {
	if err := cleanenv.ReadConfig("application.yml", c); err != nil {
		return err
	}
	return nil
}

// NewConfig initializes the configuration by reading environment variables
// and a YAML configuration file.
//
// It returns an error if there is an issue reading the environment variables
// or the configuration file.
func NewConfig() (*Config, error) {
	var cfg Config
	if err := cfg.ReadBaseConfig(); err != nil {
		return &Config{}, errors.New("NewConfg: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return &Config{}, err
	}

	appConfigName := fmt.Sprint(cfg.BaseApp.Name, "-", cfg.BaseApp.ProfileName, ".yml")

	if cfg.BaseApp.ProfileName == "dev" {
		if err := cleanenv.ReadConfig(appConfigName, &cfg); err != nil {
			return &Config{}, err
		}
	} else {
		configURL, _ := url.JoinPath(cfg.BaseApp.ConfSrvURI, appConfigName)

		data, err := common.GetFileByURL(configURL)
		if err != nil {
			return &Config{}, err
		}

		if err := cleanenv.ParseYAML(bytes.NewBuffer(data), &cfg); err != nil {
			return &Config{}, err
		}
	}

	if err := cfg.ReadPwd(); err != nil {
		return &Config{}, errors.New("Read password error: " + err.Error())
	}

	if err := cfg.validateConfig(); err != nil {
		return &Config{}, err
	}
	return &cfg, nil
}

func (c *Config) validateConfig() error {
	validate := validator.New()
	if err := validate.Struct(c); err != nil {
		return err
	}
	return nil
}

func (c *Config) ReadPwd() error {
	var pwd pwdData
	fname := fmt.Sprint("./", c.BaseApp.Name, ".json")
	if _, err := os.Stat(fname); err == nil {
		pwdFile, err := os.ReadFile(fname)
		if err != nil {
			return err
		}

		if err = json.Unmarshal(pwdFile, &pwd); err != nil {
			return err
		}

		c.App.AuthSrvPassword = pwd["app.authSrvPassword"]
	}

	return nil
}
