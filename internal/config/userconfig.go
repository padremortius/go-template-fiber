package config

import (
	"fmt"
	"os"

	gojson "github.com/goccy/go-json"
)

type (
	App struct {
		AuthSrvBaseURL  string `yaml:"authSrvBaseURL" json:"authSrvBaseURL" required:"true"`
		AuthSrvLogin    string `yaml:"authSrvLogin" json:"authSrvLogin" required:"true"`
		AuthSrvPassword string `yaml:"authSrvPassword" json:"authSrvPassword" required:"true"`
	}

	pwdData map[string]string
)

func ReadPwd() error {
	pwdFile, err := os.ReadFile(fmt.Sprint("./", Cfg.BaseApp.Name, ".json"))
	if err != nil {
		return err
	}

	if err = gojson.Unmarshal(pwdFile, &pwd); err != nil {
		return err
	}

	Cfg.App.AuthSrvPassword = pwd["app.authSrvPassword"]

	return nil
}

// func FillPubKey() error {
// 	var err error
// 	if Cfg.HTTP.Token.PublicKey, err = common.GetPubKey(Cfg.HTTP.Token.PubKeyURL); err != nil {
// 		return errors.New("Error get public key from auth-server. Error: " + err.Error())
// 	}
// 	return nil
// }
