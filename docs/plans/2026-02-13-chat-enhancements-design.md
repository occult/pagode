# Chat Enhancements: Voice Messages, Camera Capture & Security Fixes

## Overview

Extend the existing WebSocket chat with voice message recording, device camera capture for photos, and fix security issues found during review.

## Voice Messages

Record short audio clips (max 30 seconds) and send them as messages.

**Flow:** MediaRecorder API → WebM/Opus blob → upload to `/chat/upload` → URL sent via WebSocket → recipients see inline audio player.

**Frontend:**
- Mic button in ChatInput (tap to start, tap to stop)
- Recording state replaces input area with duration counter + stop button
- 30-second auto-stop with visual countdown
- Uses `navigator.mediaDevices.getUserMedia({ audio: true })`

**Backend:**
- Extend `/chat/upload` to accept `audio/webm`, `audio/ogg`, `audio/mpeg`
- Max audio file size: 2MB
- Save to `static/chat-uploads/` (same directory as images)

**Message rendering:**
- ChatMessageList detects audio URLs (extensions: `.webm`, `.ogg`, `.mp3`)
- Renders `<audio>` element with native controls inside message bubble

## Camera Capture

Take photos directly from device camera and send them in chat.

**Flow:** getUserMedia({ video }) → preview modal → canvas snapshot → blob → upload → URL message.

**Frontend:**
- Camera button in ChatInput
- Opens modal with live video preview from device camera
- Capture button takes snapshot (canvas.toBlob as JPEG)
- Uploads via same `/chat/upload` endpoint
- Closes modal and sends URL as message

**Backend:** No changes needed — reuses existing image upload.

## Security Fixes

### WebSocket Origin Check
Replace `CheckOrigin: func(r *http.Request) bool { return true }` with actual origin validation against the app's configured HTTP hostname.

### Upload Rate Limiting
Add per-IP rate limiting to `/chat/upload`: max 10 uploads per minute. Uses in-memory token bucket keyed by IP.

### File MIME Validation
Validate uploaded files by reading magic bytes (file header) instead of trusting the Content-Type header. Use Go's `http.DetectContentType()` on the first 512 bytes.

## UI Layout

```
Before: [Image] [___input___] [Send]
After:  [Image] [Camera] [Mic] [___input___] [Send]

During recording: [X Cancel] [====== 0:15 / 0:30 ======] [Stop ■]
```

## Files Changed

**Modified:**
- `pkg/handlers/chat.go` — extend upload handler, fix CORS, add rate limiting, MIME validation
- `resources/js/components/chat/ChatInput.tsx` — add camera/mic buttons, recording UI
- `resources/js/components/chat/ChatMessageList.tsx` — audio player rendering

**New:**
- `resources/js/components/chat/AudioRecorder.tsx` — recording UI component
- `resources/js/components/chat/CameraCapture.tsx` — camera preview modal
