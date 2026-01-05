package common

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

// StructToJSONBytes is ...
func StructToJSONBytes(v interface{}) ([]byte, error) {
	res, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func GetFileByURL(URL string) ([]byte, error) {
	resp, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil && err == nil {
			err = closeErr // Propagate the close error if no other error occurred
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetPubKey(URL string) (string, error) {
	rawBytes, err := GetFileByURL(URL)
	if err != nil {
		return "", err
	}
	var answer map[string]string
	if err = json.Unmarshal(rawBytes, &answer); err != nil {
		return "", err
	}
	return answer["value"], nil
}

func ReadFile(aFileName string) ([]byte, error) {
	return os.ReadFile(aFileName)
}
