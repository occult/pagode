import { useCallback, useEffect, useRef, useState } from "react";

interface UseAudioRecorderOptions {
  maxDuration?: number;
  onRecorded: (blob: Blob) => void;
}

export function useAudioRecorder({
  maxDuration = 30,
  onRecorded,
}: UseAudioRecorderOptions) {
  const [recording, setRecording] = useState(false);
  const [elapsed, setElapsed] = useState(0);
  const [starting, setStarting] = useState(false);

  const mediaRecorderRef = useRef<MediaRecorder | null>(null);
  const chunksRef = useRef<Blob[]>([]);
  const timerRef = useRef<ReturnType<typeof setInterval> | undefined>(undefined);
  const streamRef = useRef<MediaStream | null>(null);
  const onRecordedRef = useRef(onRecorded);
  onRecordedRef.current = onRecorded;

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

  useEffect(() => cleanup, [cleanup]);

  const start = useCallback(async () => {
    setStarting(true);
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true });
      streamRef.current = stream;

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
        if (e.data.size > 0) chunksRef.current.push(e.data);
      };

      recorder.onstop = () => {
        const actualType = mimeType || recorder.mimeType || "audio/webm";
        const blob = new Blob(chunksRef.current, { type: actualType });
        cleanup();
        if (blob.size > 0) onRecordedRef.current(blob);
      };

      recorder.start(100);
      setRecording(true);
      setElapsed(0);

      timerRef.current = setInterval(() => {
        setElapsed((prev) => {
          if (prev + 1 >= maxDuration) {
            recorder.stop();
            setRecording(false);
            return maxDuration;
          }
          return prev + 1;
        });
      }, 1000);
    } catch {
      alert("Microphone access denied");
    } finally {
      setStarting(false);
    }
  }, [maxDuration, cleanup]);

  const stop = useCallback(() => {
    if (mediaRecorderRef.current?.state === "recording") {
      mediaRecorderRef.current.stop();
      setRecording(false);
    }
  }, []);

  const cancel = useCallback(() => {
    if (mediaRecorderRef.current?.state === "recording") {
      mediaRecorderRef.current.ondataavailable = null;
      mediaRecorderRef.current.onstop = null;
      mediaRecorderRef.current.stop();
    }
    cleanup();
    setRecording(false);
  }, [cleanup]);

  return { recording, starting, elapsed, maxDuration, start, stop, cancel };
}
