package text

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/axpz/video-publisher/internal/analysis"
	"github.com/axpz/video-publisher/internal/config"
	genai "google.golang.org/genai"
)

type GeminiText struct {
	client *genai.Client
	model  string
}

type bearerTransport struct {
	token string
	base  http.RoundTripper
}

func (t *bearerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r := req.Clone(req.Context())
	r.Header.Set("Authorization", "Bearer "+t.token)

	base := t.base
	if base == nil {
		base = http.DefaultTransport
	}
	return base.RoundTrip(r)
}

func NewGeminiText(ctx context.Context, cfg config.Config) (*GeminiText, error) {
	var genaiCfg *genai.ClientConfig
	if cfg.LLMCfg.APIKey != "" || cfg.LLMCfg.GeminiProxyKey != "" {
		c := &genai.ClientConfig{}
		if cfg.LLMCfg.APIKey != "" {
			c.APIKey = cfg.LLMCfg.APIKey
		}
		if cfg.LLMCfg.GeminiProxyKey != "" {
			c.HTTPClient = &http.Client{
				Transport: &bearerTransport{
					token: cfg.LLMCfg.GeminiProxyKey,
				},
			}
		}
		genaiCfg = c
	}
	if cfg.LLMCfg.BaseURL != "" {
		genai.SetDefaultBaseURLs(genai.BaseURLParameters{
			GeminiURL: cfg.LLMCfg.BaseURL,
		})
	}
	client, err := genai.NewClient(ctx, genaiCfg)
	if err != nil {
		return nil, err
	}

	return &GeminiText{
		client: client,
		model:  cfg.LLMCfg.Model,
	}, nil
}

func (g *GeminiText) refine(ctx context.Context, prompt string) (*analysis.AnalysisResult, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("prompt is empty")
	}
	parts := []*genai.Part{
		{Text: prompt},
	}
	contents := []*genai.Content{
		{
			Role:  genai.RoleUser,
			Parts: parts,
		},
	}
	res, err := g.client.Models.GenerateContent(ctx, g.model, contents, nil)
	if err != nil {
		return nil, err
	}
	text := res.Text()
	if text == "" {
		return nil, errors.New("empty response from Gemini")
	}
	return parseResult(text)
}

func (g *GeminiText) RefineMetadata(ctx context.Context, description string, opts analysis.AnalyzeOptions) (*analysis.AnalysisResult, error) {
	prompt := buildPrompt(description, opts)
	return g.refine(ctx, prompt)
}

func parseResult(text string) (*analysis.AnalysisResult, error) {
	s := strings.TrimSpace(text)
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")
	if start == -1 || end == -1 || start >= end {
		return nil, fmt.Errorf("unexpected Gemini response: %q", s)
	}
	s = s[start : end+1]
	var res analysis.AnalysisResult
	if err := json.Unmarshal([]byte(s), &res); err != nil {
		return nil, err
	}
	if res.Title == "" {
		return nil, errors.New("missing title in Gemini response")
	}
	if len(res.Tags) == 0 {
		res.Tags = []string{}
	}
	return &res, nil
}
