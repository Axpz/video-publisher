package app

import (
	"github.com/axpz/video-publisher/internal/provider"
	"github.com/spf13/cobra"
)

func NewAuthCmd(factory func(platform string) (provider.VideoProvider, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "auth",
		Short: "Login to the specified platform",
		RunE: func(cmd *cobra.Command, args []string) error {
			platform, err := cmd.Root().PersistentFlags().GetString("platform")
			if err != nil {
				return err
			}
			p, err := factory(platform)
			if err != nil {
				return err
			}
			return p.Auth(cmd.Context())
		},
	}
}
