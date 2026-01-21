package text

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/axpz/video-publisher/internal/analysis"
	"github.com/axpz/video-publisher/internal/config"
)

type OpenAIText struct {
	client  *http.Client
	apiKey  string
	baseURL string
	model   string
}

func NewOpenAIText(cfg config.Config) (*OpenAIText, error) {
	baseURL := cfg.LLMCfg.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}
	model := cfg.LLMCfg.Model
	return &OpenAIText{
		client:  http.DefaultClient,
		apiKey:  cfg.LLMCfg.APIKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
	}, nil
}

func (o *OpenAIText) RefineMetadata(ctx context.Context, description string, opts analysis.AnalyzeOptions) (*analysis.AnalysisResult, error) {
	prompt := buildPrompt(description, opts)
	return o.refine(ctx, prompt)
}

func (o *OpenAIText) refine(ctx context.Context, prompt string) (*analysis.AnalysisResult, error) {
	if strings.TrimSpace(prompt) == "" {
		return nil, errors.New("prompt is empty")
	}

	body, err := json.Marshal(openAIChatRequest{
		Model: o.model,
		Messages: []openAIChatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	resp, err := o.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("openai api error: status %d: %s", resp.StatusCode, string(b))
	}

	var res openAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}
	if len(res.Choices) == 0 || res.Choices[0].Message.Content == "" {
		return nil, errors.New("empty response from OpenAI")
	}

	text := res.Choices[0].Message.Content
	return parseResult(text)
}

type openAIChatRequest struct {
	Model    string              `json:"model"`
	Messages []openAIChatMessage `json:"messages"`
}

type openAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}
