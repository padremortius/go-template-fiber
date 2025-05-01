package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type (
	App struct {
		AuthSrvPassword string `yaml:"-" json:"-" validate:"required"`
	}

	pwdData map[string]string
)

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
