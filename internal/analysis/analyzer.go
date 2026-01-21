package analysis

import "context"

type AnalyzeOptions struct {
	Language string
	Platform string
	Model    string
}

type AnalysisResult struct {
	Title    string   `json:"title"`
	Desc     string   `json:"desc"`
	Tags     []string `json:"tags"`
	Category string   `json:"category,omitempty"`
}

type VideoAnalyzer interface {
	AnalyzeVideo(ctx context.Context, videoPath string, opts AnalyzeOptions) (*AnalysisResult, error)
}

type AudioAnalyzer interface {
	AnalyzeAudio(ctx context.Context, audioPath string, opts AnalyzeOptions) (*AnalysisResult, error)
}

type TextAnalyzer interface {
	RefineMetadata(ctx context.Context, description string, opts AnalyzeOptions) (*AnalysisResult, error)
}
