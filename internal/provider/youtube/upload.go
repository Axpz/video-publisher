package youtube

import (
	"context"
	"fmt"
	"os"

	"github.com/axpz/video-publisher/internal/config"
	"google.golang.org/api/option"
	yt "google.golang.org/api/youtube/v3"
)

func (c *Client) Upload(ctx context.Context, filePath, metadataPath string) (string, error) {
	httpClient, err := c.httpClient(ctx)
	if err != nil {
		return "", err
	}

	service, err := yt.NewService(ctx, option.WithHTTPClient(httpClient))
	if err != nil {
		return "", err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return "", err
	}
	fileSize := fileInfo.Size()

	video := c.getFactoryVideoMetadata(metadataPath)

	call := service.Videos.Insert([]string{"snippet", "status", "recordingDetails"}, video)

	fmt.Printf("Start uploading file (size: %.2f MB), please wait...\n", float64(fileSize)/(1024*1024))
	fmt.Println("Using Resumable Upload, supports resuming uploads")

	response, err := call.Media(file).ProgressUpdater(func(current, total int64) {
		if total > 0 {
			percent := float64(current) * 100.0 / float64(total)
			fmt.Printf("\rUpload progress: %.2f%% (%.2f MB / %.2f MB)", percent, float64(current)/(1024*1024), float64(total)/(1024*1024))
		}
	}).Do()

	fmt.Println()

	if err != nil {
		return "", err
	}

	fmt.Printf("Upload completed successfully! Video ID: %s\n", response.Id)
	return response.Id, nil
}

func (c *Client) getFactoryVideoMetadata(metadataFile string) *yt.Video {
	var metadata, defaultMetadata struct {
		Title    string   `json:"title"`
		Desc     string   `json:"desc"`
		Tags     []string `json:"tags"`
		Category string   `json:"category"`
		Location struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"location"`
	}

	if err := config.ReadJSON(c.Config.OrigMetaFile, &defaultMetadata); err != nil {
		fmt.Printf("Failed to read default metadata file %s: %v\n", c.Config.OrigMetaFile, err)
	}

	if metadataFile != "" {
		if err := config.ReadJSON(metadataFile, &metadata); err != nil {
			fmt.Printf("Failed to read metadata file %s: %v\n", metadataFile, err)
		}
	}

	// this is by design, not a bug
	if metadata.Title == "" {
		metadata.Title = defaultMetadata.Title
	}
	if metadata.Desc == "" {
		metadata.Desc = defaultMetadata.Desc
	}

	metadata.Tags = append(metadata.Tags, defaultMetadata.Tags...)

	if metadata.Category == "" {
		metadata.Category = defaultMetadata.Category
	}

	if metadata.Location.Latitude == 0 && metadata.Location.Longitude == 0 {
		metadata.Location.Latitude = defaultMetadata.Location.Latitude
		metadata.Location.Longitude = defaultMetadata.Location.Longitude
	}

	if metadata.Title == "" {
		metadata.Title = defaultMetadata.Title
	}
	if metadata.Category == "" {
		metadata.Category = defaultMetadata.Category
	}

	return &yt.Video{
		Snippet: &yt.VideoSnippet{
			Title:           fmt.Sprintf("%s", metadata.Title),
			Description:     metadata.Desc,
			Tags:            metadata.Tags,
			CategoryId:      metadata.Category,
			DefaultLanguage: "zh-CN",
		},
		Status: &yt.VideoStatus{
			PrivacyStatus: "public",
			Embeddable:    true,
		},
		RecordingDetails: &yt.VideoRecordingDetails{
			Location: &yt.GeoPoint{
				Latitude:  metadata.Location.Latitude,
				Longitude: metadata.Location.Longitude,
			},
		},
	}
}
