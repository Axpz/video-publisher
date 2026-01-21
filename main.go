package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/axpz/video-publisher/internal/app"
	"github.com/axpz/video-publisher/internal/config"
	"github.com/axpz/video-publisher/internal/provider"
	"github.com/axpz/video-publisher/internal/provider/douyin"
	"github.com/axpz/video-publisher/internal/provider/youtube"
	"github.com/spf13/cobra"
)

var platform string

func main() {
	rootCmd := &cobra.Command{
		Use:   "video-publisher",
		Short: "Multi-platform video publisher (YouTube, Douyin)",
	}
	rootCmd.PersistentFlags().StringVarP(&platform, "platform", "p", config.DefaultPlatform, "Target platform (youtube/douyin/tiktok)")

	vpubDir := ".vpub"

	rootCmd.AddCommand(app.NewAuthCmd(func(platform string) (provider.VideoProvider, error) {
		return newProvider(platform, vpubDir)
	}))
	rootCmd.AddCommand(app.NewUploadCmd(func(platform string) (provider.VideoProvider, error) {
		return newProvider(platform, vpubDir)
	}))
	rootCmd.AddCommand(app.NewAnalyzeCmd())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("error: %v", err)
	}
}

func newProvider(platform, vpubDir string) (provider.VideoProvider, error) {
	cfg := config.Config{
		Platform:          platform,
		SessionFile:       filepath.Join(vpubDir, fmt.Sprintf("%s-session.json", platform)),
		TokenFile:         filepath.Join(vpubDir, fmt.Sprintf("%s-token.json", platform)),
		ClientSecretsFile: filepath.Join(vpubDir, fmt.Sprintf("%s-client_secrets.json", platform)),
		OrigMetaFile:      filepath.Join(vpubDir, fmt.Sprintf("%s-video_meta.orig.json", platform)),
	}

	switch platform {
	case "youtube":
		return youtube.NewClient(cfg), nil
	case "douyin":
		return douyin.NewClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", platform)
	}
}
