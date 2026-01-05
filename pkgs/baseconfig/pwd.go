package baseconfig

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/padremortius/go-template-fiber/pkgs/common"
)

func FillPwdMap(path string) (map[string]string, error) {
	pwd := make(map[string]string, 0)
	if len(path) < 1 {
		err := fmt.Errorf("path to file with secrets is empty")
		return nil, err
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return pwd, err
	}

	for _, entry := range entries {
		buff, err := common.ReadFile(filepath.Join(path, entry.Name()))
		if err != nil {
			return pwd, errors.New("Error read file: " + entry.Name() + " Error: " + err.Error())
		}
		pwd[entry.Name()] = strings.Trim(string(buff), "\n")
	}
	return pwd, nil
}
