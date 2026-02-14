import { useCallback, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ImagePlus, Camera, Send, Loader2, Mic, X, Square } from "lucide-react";
import { useAudioRecorder } from "@/hooks/useAudioRecorder";
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

function formatTime(seconds: number) {
  const m = Math.floor(seconds / 60);
  const s = seconds % 60;
  return `${m}:${s.toString().padStart(2, "0")}`;
}

export function ChatInput({ onSend, onTyping, disabled }: ChatInputProps) {
  const [value, setValue] = useState("");
  const [uploading, setUploading] = useState(false);
  const [cameraOpen, setCameraOpen] = useState(false);
  const typingTimeout = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const fileInputRef = useRef<HTMLInputElement>(null);

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

  const {
    recording,
    starting,
    elapsed,
    maxDuration,
    start: startRecording,
    stop: stopRecording,
    cancel: cancelRecording,
  } = useAudioRecorder({ maxDuration: 30, onRecorded: handleAudioRecorded });

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
      if (!typingTimeout.current) onTyping();
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
  const pct = maxDuration > 0 ? (elapsed / maxDuration) * 100 : 0;

  return (
    <>
      {/* Relative wrapper — the recording overlay is positioned inside this */}
      <div className="relative border-t">
        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png,image/gif,image/webp"
          className="hidden"
          onChange={handleFileSelect}
        />

        {/* Base input bar — always in flow, dimensions never change */}
        <form
          onSubmit={handleSubmit}
          className="flex items-center gap-1.5 sm:gap-2 px-2 sm:px-3 py-2"
        >
          <Button
            type="button"
            size="icon"
            variant="ghost"
            disabled={isDisabled}
            onClick={() => fileInputRef.current?.click()}
            title="Send image"
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9"
          >
            {uploading ? (
              <Loader2 className="h-5 w-5 sm:h-4 sm:w-4 animate-spin" />
            ) : (
              <ImagePlus className="h-5 w-5 sm:h-4 sm:w-4" />
            )}
          </Button>

          <Button
            type="button"
            size="icon"
            variant="ghost"
            disabled={isDisabled}
            onClick={() => setCameraOpen(true)}
            title="Take photo"
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9"
          >
            <Camera className="h-5 w-5 sm:h-4 sm:w-4" />
          </Button>

          <Button
            type="button"
            size="icon"
            variant="ghost"
            disabled={isDisabled || starting}
            onClick={startRecording}
            title="Record voice message"
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9"
          >
            {starting ? (
              <Loader2 className="h-5 w-5 sm:h-4 sm:w-4 animate-spin" />
            ) : (
              <Mic className="h-5 w-5 sm:h-4 sm:w-4" />
            )}
          </Button>

          <Input
            value={value}
            onChange={handleChange}
            placeholder="Type a message..."
            disabled={isDisabled}
            autoComplete="off"
            className="flex-1 min-w-0 h-10 sm:h-9"
          />

          <Button
            type="submit"
            size="icon"
            disabled={isDisabled || !value.trim()}
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9"
          >
            <Send className="h-5 w-5 sm:h-4 sm:w-4" />
          </Button>
        </form>

        {/* Recording overlay — absolutely positioned, ZERO layout impact on the form */}
        <div
          className={`absolute inset-0 z-10 flex items-center gap-2 px-2 sm:px-3 py-2 bg-background transition-opacity duration-200 ${
            recording ? "opacity-100" : "opacity-0 pointer-events-none"
          }`}
        >
          <Button
            type="button"
            size="icon"
            variant="ghost"
            onClick={cancelRecording}
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9"
          >
            <X className="h-5 w-5 sm:h-4 sm:w-4" />
          </Button>

          <div className="flex-1 flex items-center gap-2 min-w-0">
            <span className="h-2.5 w-2.5 rounded-full bg-red-500 animate-pulse flex-shrink-0" />
            <span className="text-xs font-mono text-muted-foreground whitespace-nowrap">
              {formatTime(elapsed)}/{formatTime(maxDuration)}
            </span>
            <div className="flex-1 bg-muted rounded-full h-1.5 min-w-0">
              <div
                className="bg-red-500 h-1.5 rounded-full transition-all duration-1000"
                style={{ width: `${pct}%` }}
              />
            </div>
          </div>

          <Button
            type="button"
            size="icon"
            variant="destructive"
            onClick={stopRecording}
            className="flex-shrink-0 h-10 w-10 sm:h-9 sm:w-9 rounded-full"
          >
            <Square className="h-4 w-4 sm:h-3 sm:w-3" />
          </Button>
        </div>
      </div>

      <CameraCapture
        open={cameraOpen}
        onCapture={handleCameraCapture}
        onClose={() => setCameraOpen(false)}
      />
    </>
  );
}
