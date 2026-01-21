package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Platform          string
	SessionFile       string
	TokenFile         string
	ClientSecretsFile string

	OrigMetaFile string
	LLMCfg       LLMConfig
}

type LLMConfig struct {
	BaseURL        string `json:"base_url"`
	Lang           string `json:"lang"`
	Model          string `json:"model"`
	APIKey         string `json:"api_key"`
	GeminiProxyKey string `json:"gemini_proxy_key"`
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

func LoadLLMConfig(path string) (LLMConfig, error) {
	var cfg LLMConfig
	if err := ReadJSON(path, &cfg); err != nil {
		return LLMConfig{}, err
	}
	return cfg, nil
}
