package app

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/axpz/video-publisher/internal/analysis"
	"github.com/axpz/video-publisher/internal/analysis/text"
	"github.com/axpz/video-publisher/internal/config"
	"github.com/spf13/cobra"
)

func NewAnalyzeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "analyze [description]",
		Short: "Generate video metadata JSON from a short description using Gemini",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			description := strings.Join(args, " ")
			ctx := cmd.Context()
			platform, err := cmd.Root().PersistentFlags().GetString("platform")
			if err != nil {
				return err
			}
			llmPath := filepath.Join(".vpub", "llm.json")
			llmCfg, err := config.LoadLLMConfig(llmPath)
			if err != nil {
				return err
			}
			cfg := config.Config{
				Platform: platform,
				LLMCfg:   llmCfg,
			}
			return runAnalyze(ctx, cfg, description)
		},
	}

	return cmd
}

func runAnalyze(ctx context.Context, cfg config.Config, description string) error {
	now := time.Now().Format("20060102150405")
	outPath := fmt.Sprintf("./%s-meta.%s.json", cfg.Platform, now)

	lang := cfg.LLMCfg.Lang
	analyzer, err := newTextAnalyzer(ctx, cfg)
	if err != nil {
		return err
	}
	meta, err := analyzer.RefineMetadata(ctx, description, analysis.AnalyzeOptions{
		Language: lang,
		Platform: cfg.Platform,
		Model:    cfg.LLMCfg.Model,
	})
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Printf("Result decoded into %s", outPath)
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(meta)
}

func newTextAnalyzer(ctx context.Context, cfg config.Config) (analysis.TextAnalyzer, error) {
	model := cfg.LLMCfg.Model
	if model == "" {
		log.Fatal("llm model is empty")
	}

	lower := strings.ToLower(model)
	if strings.Contains(lower, "gpt") || strings.Contains(lower, "o3") {
		return text.NewOpenAIText(cfg)
	}
	return text.NewGeminiText(ctx, cfg)
}
