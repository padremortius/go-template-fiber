package baseconfig

import (
	"fmt"
)

type (
	BaseApp struct {
		Name        string `env-required:"true" yaml:"name" env:"APP_NAME" json:"name" env-update:"true"`
		ProfileName string `env-required:"true" yaml:"profileName" json:"profileName" env:"PROFILE_NAME" env-update:"true"`
		ConfSrvURI  string `yaml:"confSrvURI" json:"confSrvURI" env:"CONF_SRV_URI" env-update:"true"`
		SecPath     string `yaml:"secPath" json:"secPath" env:"SEC_PATH" env-update:"true"`
	}

	Version struct {
		BuildVersion   string `json:"buildVersion"`
		BuildTimeStamp string `json:"buildTimeStamp,omitempty"`
		GitBranch      string `json:"gitBranch,omitempty"`
		GitHash        string `json:"gitHash,omitempty"`
	}
)

var (
	binVersion = "0.0.1"
)

func InitVersion(aBuildNumber, aBuildTimeStamp, aGitBranch, aGitHash string) *Version {
	return &Version{
		BuildVersion:   fmt.Sprint(binVersion, ".", aBuildNumber),
		GitBranch:      aGitBranch,
		GitHash:        aGitHash,
		BuildTimeStamp: aBuildTimeStamp,
	}
}
