package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	BaseApp struct {
		Name        string `env-required:"true" yaml:"name" env:"APP_NAME" json:"name" env-update:"true"`
		ProfileName string `env-required:"true" yaml:"profileName" json:"profileName" env:"PROFILE_NAME" env-update:"true"`
		ConfSrvURI  string `yaml:"confSrvURI" json:"confSrvURI" env:"CONF_SRV_URI" env-update:"true"`
	}

	Version struct {
		Version        string `json:"version"`
		BuildTimeStamp string `json:"buildTimeStamp,omitempty"`
		GitBranch      string `json:"gitBranch,omitempty"`
		GitHash        string `json:"gitHash,omitempty"`
	}
)

var (
	binVersion = "0.0.1"
)

func (c *Config) ReadBaseConfig() error {
	if err := cleanenv.ReadConfig("application.yml", c); err != nil {
		return err
	}
	return nil
}

func InitVersion(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash string) *Version {
	return &Version{
		Version:        fmt.Sprint(binVersion, ".", aBuildNumber),
		GitBranch:      aGitBranch,
		GitHash:        aGitHash,
		BuildTimeStamp: aBuildTimeStamp,
	}
}
