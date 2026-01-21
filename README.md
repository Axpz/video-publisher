# video-publisher

A command-line tool for one-click video publishing to multiple platforms (YouTube, Douyin, etc.) with AI-powered metadata generation.

## Features

- **Multi-platform support**: `youtube`, `douyin` (TikTok support coming soon).
- **AI Analysis**: Generate video titles, descriptions, and tags from a simple sentence using Gemini or OpenAI.
- **Unified Workflow**: Simple `auth` -> `analyze` -> `upload` pipeline.
- **Configurable**: Flexible metadata and platform-specific settings.
- **Go-powered**: Single binary, high performance, and easy to deploy.

## Installation

```bash
make build
```

This will generate the `video-publisher` binary in the root directory.

## Usage

### 1. Authentication

First, authenticate with your target platform:

```bash
# Default is YouTube
./video-publisher auth

# For Douyin
./video-publisher auth -p douyin
```

Credentials and session data will be stored in the `.vpub/` directory.

### 2. AI Metadata Analysis (Optional)

Instead of manually writing JSON metadata, you can generate it using LLM:

```bash
./video-publisher analyze "一个关于在山东旅行并品尝当地美食的短视频"
```

This command will:
1. Read LLM configuration from `.vpub/llm_secrets.json`.
2. Use Gemini (default) or OpenAI to generate a optimized title, description, and tags.
3. Save the result to a file like `youtube-meta.20260121120000.json`.

### 3. Upload Video

Upload your video using the generated or manual metadata:

```bash
# Format: ./video-publisher upload <video_path> <metadata_path>
./video-publisher upload my-trip.mp4 youtube-meta.20260121120000.json
```

## Configuration

The tool stores all configuration and secrets in the `.vpub/` directory.

### LLM Configuration (`.vpub/llm.json`)

To use the `analyze` command, create this file:

```json
{
  "model": "gemini-2.0-flash-exp",
  "api_key": "YOUR_GEMINI_API_KEY",
  "base_url": "",
  "lang": "zh-CN",
  "gemini_proxy_key": ""
}
```

- **model**: Supports Gemini models (e.g., `gemini-2.0-flash-exp`) or OpenAI models (e.g., `gpt-4o`, `o3-mini`).
- **api_key**: Your provider API key.
- **base_url**: Optional. Custom API endpoint (useful for proxies).
- **lang**: Target language for generated metadata.
- **gemini_proxy_key**: Optional. Specific for certain proxy setups.

### Metadata Format

The metadata file (passed to `upload`) follows this structure:

```json
{
  "title": "视频标题",
  "desc": "视频详细描述",
  "tags": ["标签1", "标签2"],
  "category": "28"
}
```

*Note: `category` is platform-specific (e.g., YouTube category IDs).*

## Roadmap

- [x] AI auto-generate metadata for SEO
- [ ] TikTok platform support
- [ ] Batch upload support
- [ ] AI video editing and auto-captioning
