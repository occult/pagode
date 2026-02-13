import { useCallback, useEffect, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Mic, Square, X, Loader2 } from "lucide-react";

interface AudioRecorderProps {
  onRecorded: (blob: Blob) => void;
  onCancel: () => void;
  onRecordingChange?: (recording: boolean) => void;
  maxDuration?: number;
  disabled?: boolean;
}

export function AudioRecorder({ onRecorded, onCancel, onRecordingChange, maxDuration = 30, disabled }: AudioRecorderProps) {
  const [recording, setRecording] = useState(false);
  const [elapsed, setElapsed] = useState(0);
  const [starting, setStarting] = useState(false);
  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<ReturnType<typeof setInterval> | undefined>(undefined);
  const streamRef = useRef<MediaStream | null>(null);

  // Use refs for callbacks to avoid stale closures in MediaRecorder.onstop
  const onRecordedRef = useRef(onRecorded);
  onRecordedRef.current = onRecorded;
  const onRecordingChangeRef = useRef(onRecordingChange);
  onRecordingChangeRef.current = onRecordingChange;

  const cleanup = useCallback(() => {
    if (timerRef.current) {
      clearInterval(timerRef.current);
      timerRef.current = undefined;
    }
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((t) => t.stop());
      streamRef.current = null;
    }
    mediaRecorderRef.current = null;
    chunksRef.current = [];
  }, []);

  useEffect(() => {
    return cleanup;
  }, [cleanup]);

  const startRecording = useCallback(async () => {
    setStarting(true);
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      streamRef.current = stream;

      // Safari/iOS uses mp4, Chrome/Firefox use webm
      const mimeType = MediaRecorder.isTypeSupported("audio/webm;codecs=opus")
        ? "audio/webm;codecs=opus"
        : MediaRecorder.isTypeSupported("audio/webm")
          ? "audio/webm"
          : MediaRecorder.isTypeSupported("audio/mp4")
            ? "audio/mp4"
            : "";

      const recorder = mimeType
        ? new MediaRecorder(stream, { mimeType })
        : new MediaRecorder(stream);
      mediaRecorderRef.current = recorder;
      chunksRef.current = [];

      recorder.ondataavailable = (e) => {
        if (e.data.size > 0) {
          chunksRef.current.push(e.data);
        }
      };

      recorder.onstop = () => {
        const actualType = mimeType || recorder.mimeType || "audio/webm";
        const blob = new Blob(chunksRef.current, { type: actualType });
        cleanup();
        if (blob.size > 0) {
          // Use ref to always call the latest callback
          onRecordedRef.current(blob);
        }
      };

      recorder.start(100);
      setRecording(true);
      setElapsed(0);
      onRecordingChangeRef.current?.(true);

      timerRef.current = setInterval(() => {
        setElapsed((prev) => {
          if (prev + 1 >= maxDuration) {
            recorder.stop();
            setRecording(false);
            onRecordingChangeRef.current?.(false);
            return maxDuration;
          }
          return prev + 1;
        });
      }, 1000);
    } catch {
      alert("Microphone access denied");
      onCancel();
    } finally {
      setStarting(false);
    }
  }, [maxDuration, onCancel, cleanup]);

  const stopRecording = useCallback(() => {
    if (mediaRecorderRef.current && mediaRecorderRef.current.state === "recording") {
      mediaRecorderRef.current.stop();
      setRecording(false);
      onRecordingChangeRef.current?.(false);
    }
  }, []);

  const handleCancel = useCallback(() => {
    if (mediaRecorderRef.current && mediaRecorderRef.current.state === "recording") {
      mediaRecorderRef.current.ondataavailable = null;
      mediaRecorderRef.current.onstop = null;
      mediaRecorderRef.current.stop();
    }
    cleanup();
    setRecording(false);
    onRecordingChangeRef.current?.(false);
    onCancel();
  }, [cleanup, onCancel]);

  const formatTime = (seconds: number) => {
    const m = Math.floor(seconds / 60);
    const s = seconds % 60;
    return `${m}:${s.toString().padStart(2, "0")}`;
  };

  if (!recording && !starting) {
    return (
      <Button
        type="button"
        size="icon"
        variant="ghost"
        disabled={disabled}
        onClick={startRecording}
        title="Record voice message"
      >
        <Mic className="h-4 w-4" />
      </Button>
    );
  }

  if (starting) {
    return (
      <Button type="button" size="icon" variant="ghost" disabled>
        <Loader2 className="h-4 w-4 animate-spin" />
      </Button>
    );
  }

  return (
    <div className="flex items-center gap-2 flex-1">
      <Button type="button" size="icon" variant="ghost" onClick={handleCancel}>
        <X className="h-4 w-4" />
      </Button>
      <div className="flex-1 flex items-center gap-2">
        <span className="h-2 w-2 rounded-full bg-red-500 animate-pulse" />
        <span className="text-sm font-mono text-muted-foreground">
          {formatTime(elapsed)} / {formatTime(maxDuration)}
        </span>
        <div className="flex-1 bg-muted rounded-full h-1.5">
          <div
            className="bg-red-500 h-1.5 rounded-full transition-all"
            style={{ width: `${(elapsed / maxDuration) * 100}%` }}
          />
        </div>
      </div>
      <Button type="button" size="icon" variant="destructive" onClick={stopRecording}>
        <Square className="h-3 w-3" />
      </Button>
    </div>
  );
}
