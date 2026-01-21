package app

import (
	"github.com/axpz/video-publisher/internal/provider"
	"github.com/spf13/cobra"
)

func NewUploadCmd(factory func(platform string) (provider.VideoProvider, error)) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upload [video_file] [metadata_file]",
		Short: "Upload video to the specified platform, i.e upload video.mp4 metadata_file.json",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			platform, err := cmd.Root().PersistentFlags().GetString("platform")
			if err != nil {
				return err
			}
			p, err := factory(platform)
			if err != nil {
				return err
			}
			_, err = p.Upload(cmd.Context(), args[0], args[1])
			return err
		},
	}
	return cmd
}
