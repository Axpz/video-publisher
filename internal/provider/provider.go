package provider

import "context"

// VideoProvider defines the standard behavior for all video providers.
type VideoProvider interface {
	// Auth handles the login or authentication process for the platform.
	Auth(ctx context.Context) error

	// Upload uploads a local video file to the platform.
	// It returns the unique identifier or URL of the uploaded video,
	// along with any potential errors.
	Upload(ctx context.Context, filePath, metadataPath string) (string, error)
}
