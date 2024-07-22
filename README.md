# mediaconv

For now, a simple direct URL to GIF conversion API.

## Run

1. Copy and configure `config.example.json`
2. Build with `go build cmd/mediaconv/main.go`
3. Run with `./main` (or `.\main.exe` on Windows)

## Usage

`POST http://localhost:$port/`
```json
{
  "url": "https://link.to.websi.te/with/video.mp4"
}
```
