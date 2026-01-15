package app

import (
	"encoding/json"
	"os"
)

type Config struct {
	DefaultMetadataFile string

	SessionFile       string
	TokenFile         string
	ClientSecretsFile string
}

var DefaultPlatform = "youtube"

func ReadJSON(path string, v any) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(v)
}
