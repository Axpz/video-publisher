package main

import (
	"fmt"
	"os"

	"github.com/axpz/video-publisher/internal/app"
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

	rootCmd.PersistentFlags().StringVarP(&platform, "platform", "p", app.DefaultPlatform, "Target platform (youtube/douyin/tiktok)")

	rootCmd.AddCommand(newAuthCmd())
	rootCmd.AddCommand(newUploadCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newAuthCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Login to the specified platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newVideoProvider(platform)
			if err != nil {
				return err
			}
			fmt.Printf("Start authentication for %s...\n", platform)
			return p.Auth()
		},
	}
}

func newUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload [video_file] [metadata_file]",
		Short: "Upload video to the specified platform, i.e upload video.mp4 metadata_file.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			p, err := newVideoProvider(platform)
			if err != nil {
				return err
			}
			fmt.Printf("Start upload to %s...\n", platform)
			_, err = p.Upload(args[0], args[1])
			return err
		},
	}
}

func newVideoProvider(name string) (provider.VideoProvider, error) {
	cfg := app.Config{
		SessionFile:       fmt.Sprintf(".auth/%s_session.json", name),
		TokenFile:         fmt.Sprintf(".auth/%s_token.json", name),
		ClientSecretsFile: fmt.Sprintf(".auth/%s_client_secrets.json", name),
	}

	switch name {
	case "youtube":
		cfg.DefaultMetadataFile = "metadata/youtube_metadata.json"
		return youtube.NewClient(cfg), nil
	case "douyin":
		return douyin.NewClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported platform: %s", name)
	}
}
