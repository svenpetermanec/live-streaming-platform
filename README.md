
This project is a minimal live streaming platform built to ingest real-time video via SRT, transcode it to multiple HLS formats, and serve it to viewers. It’s composed of two core services:

- `transcoder`: Ingests video streams via SRT, transcodes to HLS.
- `api`: Handles stream creation and serves HLS content to clients.

---
![Blank diagram(2)](https://github.com/user-attachments/assets/d39b7a6c-c6e6-476b-ac3c-931434ca15ad)


## Example Flow

1. **Streamer** hits `/account/create?username=someguy`
2. Receives `streamId=xyz123`
3. Begins streaming via SRT to the `transcoder`, tagging stream with `xyz123`
4. `transcoder` fetches `someguy` from Redis and begins `ffmpeg` transcoding
5. Files saved locally in HLS format
6. **Clients** view via `/streams/someguy/.../720p30.m3u8` 

---

## Getting Started

1. **Run services and infra**
 ```bash
docker compose up --build
  ```

2. Create stream
  ```bash
curl -X POST "http://localhost:8080/account/create?username=someguy"
  ```

3. Start streaming
Stream by using any SRT-compatible encoder. Here’s an example using FFmpeg in Docker:
  ```bash
docker run --rm --network srt_network -v $(pwd):/config linuxserver/ffmpeg:latest -re -i /sample/video.mp4 -c copy -f mpegts "srt://transcoder:5270?streamid=xyz123"
  ```

4. Watch the stream
 Watch using a browser with native HLS support (e.g. Safari), a player using hls.js, VLC media player...
```bash
http://localhost:8080/streams/someguy/{streamTimestamp}/{quality}/playlist.m3u8
```
