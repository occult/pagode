import { useCallback, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ImagePlus, Camera, Send, Loader2 } from "lucide-react";
import { AudioRecorder } from "./AudioRecorder";
import { CameraCapture } from "./CameraCapture";

interface ChatInputProps {
  onSend: (message: string) => void;
  onTyping: () => void;
  disabled?: boolean;
}

function getCsrfToken(): string | undefined {
  return document.cookie
    .split("; ")
    .find((row) => row.startsWith("XSRF-TOKEN="))
    ?.split("=")[1];
}

async function uploadFile(blob: Blob, filename: string): Promise<string> {
  const formData = new FormData();
  formData.append("file", blob, filename);

  const csrfToken = getCsrfToken();
  const res = await fetch("/chat/upload", {
    method: "POST",
    body: formData,
    headers: csrfToken ? { "X-XSRF-TOKEN": decodeURIComponent(csrfToken) } : {},
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || "Upload failed");
  }

  const data = await res.json();
  return data.url;
}

export function ChatInput({ onSend, onTyping, disabled }: ChatInputProps) {
  const [value, setValue] = useState("");
  const [uploading, setUploading] = useState(false);
  const [recording, setRecording] = useState(false);
  const [cameraOpen, setCameraOpen] = useState(false);
  const typingTimeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      const trimmed = value.trim();
      if (!trimmed) return;
      onSend(trimmed);
      setValue("");
    },
    [value, onSend]
  );

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setValue(e.target.value);
      if (!typingTimeout.current) {
        onTyping();
      }
      clearTimeout(typingTimeout.current);
      typingTimeout.current = setTimeout(() => {
        typingTimeout.current = undefined;
      }, 2000);
    },
    [onTyping]
  );

  const handleFileSelect = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      const file = e.target.files?.[0];
      if (!file) return;
      e.target.value = "";

      if (file.size > 5 * 1024 * 1024) {
        alert("File too large (max 5MB)");
        return;
      }

      setUploading(true);
      try {
        const url = await uploadFile(file, file.name);
        onSend(url);
      } catch (err) {
        console.error("Upload failed:", err);
        alert("Failed to upload file");
      } finally {
        setUploading(false);
      }
    },
    [onSend]
  );

  const handleAudioRecorded = useCallback(
    async (blob: Blob) => {
      setUploading(true);
      try {
        const ext = blob.type.includes("mp4") ? ".m4a" : blob.type.includes("ogg") ? ".ogg" : ".webm";
        const url = await uploadFile(blob, `voice-message${ext}`);
        onSend(url);
      } catch (err) {
        console.error("Audio upload failed:", err);
        alert("Failed to upload voice message");
      } finally {
        setUploading(false);
      }
    },
    [onSend]
  );

  const handleCameraCapture = useCallback(
    async (blob: Blob) => {
      setCameraOpen(false);
      setUploading(true);
      try {
        const url = await uploadFile(blob, "camera-photo.jpg");
        onSend(url);
      } catch (err) {
        console.error("Camera upload failed:", err);
        alert("Failed to upload photo");
      } finally {
        setUploading(false);
      }
    },
    [onSend]
  );

  const isDisabled = disabled || uploading;

  return (
    <>
      <form onSubmit={handleSubmit} className="flex items-center gap-2 p-4 border-t">
        {/* Hidden file input for image picker */}
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          className="hidden"
          onChange={handleFileSelect}
        />

        {/* Uploading indicator */}
        {uploading && !recording && (
          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
        )}

        {/* Image + Camera buttons: hidden when recording or uploading */}
        {!recording && !uploading && (
          <>
            <Button
              type="button"
              size="icon"
              variant="ghost"
              disabled={isDisabled}
              onClick={() => fileInputRef.current?.click()}
              title="Send image"
            >
              <ImagePlus className="h-4 w-4" />
            </Button>

            <Button
              type="button"
              size="icon"
              variant="ghost"
              disabled={isDisabled}
              onClick={() => setCameraOpen(true)}
              title="Take photo"
            >
              <Camera className="h-4 w-4" />
            </Button>
          </>
        )}

        {/* Single AudioRecorder instance — always mounted, manages its own UI */}
        <AudioRecorder
          onRecorded={handleAudioRecorded}
          onCancel={() => {}}
          onRecordingChange={setRecording}
          maxDuration={30}
          disabled={isDisabled}
        />

        {/* Text input + Send: hidden when recording */}
        {!recording && (
          <>
            <Input
              value={value}
              onChange={handleChange}
              placeholder="Type a message..."
              disabled={isDisabled}
              autoComplete="off"
              className="flex-1"
            />
            <Button type="submit" size="icon" disabled={isDisabled || !value.trim()}>
              <Send className="h-4 w-4" />
            </Button>
          </>
        )}
      </form>

      <CameraCapture
        open={cameraOpen}
        onCapture={handleCameraCapture}
        onClose={() => setCameraOpen(false)}
      />
    </>
  );
}
