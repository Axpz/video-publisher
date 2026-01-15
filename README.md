# video-publisher

A command-line tool for one-click video publishing to multiple platforms YouTube/Douyin/...

## Features

- Multi-platform support: `youtube`, `douyin`
- Unified `auth` and `upload` commands
- Configurable metadata via JSON (title, description, etc.)
- Written in Go, single binary executable

## Quick Start

## Usage

### 1. Login

For YouTube (default platform):

```bash
./video-publisher auth
```

Specify platform (e.g., Douyin):

```bash
./video-publisher auth -p douyin
```

The tool uses platform-specific authentication files stored in:

- Session file: `.auth/<platform>_session.json`
- Token file: `.auth/<platform>_token.json`
- Client secrets: `.auth/<platform>_client_secrets.json`

Where `<platform>` is `youtube` or `douyin`.

Running `auth` generates/updates these files, which are automatically used by `upload`.

### 2. Upload Video

Format:

```bash
./video-publisher upload <video_path> <metadata_path>
```

Example (YouTube):

```bash
./video-publisher upload video.mp4 video-meta.json
```

## Roadmap

- [ ] AI auto-generate metadata for SEO (In Development)
- [ ] AI video editing
