package common

import (
	"io"
	"net/http"

	gojson "github.com/goccy/go-json"
)

// StructToJSONBytes is ...
func StructToJSONBytes(v interface{}) ([]byte, error) {
	res, err := gojson.Marshal(v)
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
	defer resp.Body.Close()

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
	if err = gojson.Unmarshal(rawBytes, &answer); err != nil {
		return "", err
	}
	return answer["value"], nil
}
