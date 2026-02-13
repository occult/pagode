import { useCallback, useEffect, useRef, useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Camera, SwitchCamera, X } from "lucide-react";

interface CameraCaptureProps {
  open: boolean;
  onCapture: (blob: Blob) => void;
  onClose: () => void;
}

export function CameraCapture({ open, onCapture, onClose }: CameraCaptureProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const streamRef = useRef<MediaStream | null>(null);
  const [facingMode, setFacingMode] = useState<"user" | "environment">("environment");
  const [ready, setReady] = useState(false);

  const stopStream = useCallback(() => {
    if (streamRef.current) {
      streamRef.current.getTracks().forEach((t) => t.stop());
      streamRef.current = null;
    }
  }, []);

  const startCamera = useCallback(async (facing: "user" | "environment") => {
    stopStream();
    setReady(false);

    try {
      // Try with facingMode first, fall back to any camera if it fails
      let stream: MediaStream;
      try {
        stream = await navigator.mediaDevices.getUserMedia({
          video: { facingMode: { ideal: facing }, width: { ideal: 1280 }, height: { ideal: 720 } },
          audio: false,
        });
      } catch {
        // Fallback: request any available camera
        stream = await navigator.mediaDevices.getUserMedia({ video: true, audio: false });
      }

      streamRef.current = stream;

      const video = videoRef.current;
      if (video) {
        video.srcObject = stream;
        // iOS Safari requires explicit play() after srcObject assignment
        video.onloadedmetadata = () => {
          video.play().then(() => setReady(true)).catch(() => setReady(true));
        };
      }
    } catch {
      alert("Camera access denied");
      onClose();
    }
  }, [onClose, stopStream]);

  useEffect(() => {
    if (open) {
      // Small delay to let Dialog render the video element first
      const timeout = setTimeout(() => startCamera(facingMode), 100);
      return () => {
        clearTimeout(timeout);
        stopStream();
        setReady(false);
      };
    }
    stopStream();
    setReady(false);
  }, [open, facingMode, startCamera, stopStream]);

  const handleCapture = useCallback(() => {
    const video = videoRef.current;
    const canvas = canvasRef.current;
    if (!video || !canvas) return;

    canvas.width = video.videoWidth;
    canvas.height = video.videoHeight;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    ctx.drawImage(video, 0, 0);
    canvas.toBlob(
      (blob) => {
        if (blob) {
          stopStream();
          onCapture(blob);
        }
      },
      "image/jpeg",
      0.85
    );
  }, [onCapture, stopStream]);

  const toggleCamera = useCallback(() => {
    setFacingMode((prev) => (prev === "user" ? "environment" : "user"));
  }, []);

  return (
    <Dialog open={open} onOpenChange={(v) => !v && onClose()}>
      <DialogContent className="max-w-[calc(100vw-2rem)] sm:max-w-2xl h-[calc(100dvh-2rem)] sm:h-auto sm:max-h-[90dvh] p-0 overflow-hidden flex flex-col">
        <DialogHeader className="p-4 pb-0">
          <DialogTitle>Take a Photo</DialogTitle>
        </DialogHeader>
        <div className="relative bg-black flex-1 min-h-0">
          <video
            ref={videoRef}
            autoPlay
            playsInline
            muted
            className="w-full h-full object-cover sm:aspect-video sm:h-auto"
            style={{ WebkitTransform: "translateZ(0)" }}
          />
          <canvas ref={canvasRef} className="hidden" />
          {!ready && (
            <div className="absolute inset-0 flex items-center justify-center text-white text-sm">
              Starting camera...
            </div>
          )}
        </div>
        <div className="flex items-center justify-center gap-4 p-4">
          <Button type="button" size="icon" variant="outline" onClick={onClose}>
            <X className="h-4 w-4" />
          </Button>
          <Button
            type="button"
            size="lg"
            className="rounded-full h-14 w-14"
            onClick={handleCapture}
            disabled={!ready}
          >
            <Camera className="h-6 w-6" />
          </Button>
          <Button type="button" size="icon" variant="outline" onClick={toggleCamera}>
            <SwitchCamera className="h-4 w-4" />
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
