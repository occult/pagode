import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";
import { Loader2, Play, Pause } from "lucide-react";
import { getAvatarColor, getInitials } from "./avatar-colors";

interface MessageItem {
  type: "message" | "join" | "leave";
  id?: number;
  senderName?: string;
  body?: string;
  createdAt?: string;
}

interface ChatMessageListProps {
  messages: MessageItem[];
  currentUser: string;
  hasMore: boolean;
  loadingMore: boolean;
  onLoadMore: () => void;
}

function formatTime(dateStr?: string): string {
  if (!dateStr) return "";
  const d = new Date(dateStr);
  return d.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" });
}

const audioMimeMap: Record<string, string> = {
  ".webm": "audio/webm",
  ".ogg": "audio/ogg",
  ".mp3": "audio/mpeg",
  ".m4a": "audio/mp4",
  ".aac": "audio/aac",
};

function getMediaInfo(text?: string): { type: "image"; src: string } | { type: "audio"; src: string; mime: string } | null {
  if (!text || !text.startsWith("/files/chat-uploads/")) return null;
  const lower = text.toLowerCase();

  for (const [ext, mime] of Object.entries(audioMimeMap)) {
    if (lower.endsWith(ext)) {
      return { type: "audio", src: text, mime };
    }
  }

  if (/\.(jpe?g|png|gif|webp)$/.test(lower)) {
    return { type: "image", src: text };
  }

  return null;
}

// Custom audio player with progress bar
function AudioPlayer({ src, mime, isMine }: { src: string; mime: string; isMine: boolean }) {
  const audioRef = useRef<HTMLAudioElement>(null);
  const [playing, setPlaying] = useState(false);
  const [progress, setProgress] = useState(0);
  const [duration, setDuration] = useState(0);
  const rafRef = useRef<number>(0);
  const resolvedDurationRef = useRef(false);

  const updateProgress = useCallback(() => {
    const audio = audioRef.current;
    if (audio && playing) {
      setProgress(audio.currentTime);
      rafRef.current = requestAnimationFrame(updateProgress);
    }
  }, [playing]);

  useEffect(() => {
    if (playing) {
      rafRef.current = requestAnimationFrame(updateProgress);
    }
    return () => cancelAnimationFrame(rafRef.current);
  }, [playing, updateProgress]);

  // MediaRecorder files (WebM/M4A) often lack duration in the header.
  // Force the browser to resolve it by briefly seeking to the end.
  useEffect(() => {
    const audio = audioRef.current;
    if (!audio) return;

    const tryResolveDuration = () => {
      if (resolvedDurationRef.current) return;
      const d = audio.duration;
      if (d && isFinite(d) && d > 0) {
        resolvedDurationRef.current = true;
        setDuration(d);
        return;
      }
      // Duration is Infinity or 0 — seek to a large value to force the browser
      // to determine the real length, then seek back to 0.
      audio.currentTime = 1e6;
    };

    const onSeeked = () => {
      if (resolvedDurationRef.current) return;
      const d = audio.duration;
      if (d && isFinite(d) && d > 0) {
        resolvedDurationRef.current = true;
        setDuration(d);
      }
      // Reset to the beginning (only if not playing)
      if (!audio.paused) return;
      audio.currentTime = 0;
    };

    const onDurationChange = () => {
      const d = audio.duration;
      if (d && isFinite(d) && d > 0) {
        resolvedDurationRef.current = true;
        setDuration(d);
      }
    };

    audio.addEventListener("loadedmetadata", tryResolveDuration);
    audio.addEventListener("durationchange", onDurationChange);
    audio.addEventListener("seeked", onSeeked);

    // If already loaded (cached), try immediately
    if (audio.readyState >= 1) tryResolveDuration();

    return () => {
      audio.removeEventListener("loadedmetadata", tryResolveDuration);
      audio.removeEventListener("durationchange", onDurationChange);
      audio.removeEventListener("seeked", onSeeked);
    };
  }, [src]);

  const togglePlay = () => {
    const audio = audioRef.current;
    if (!audio) return;
    if (playing) {
      audio.pause();
      setPlaying(false);
    } else {
      audio.play();
      setPlaying(true);
    }
  };

  const handleEnded = () => {
    setPlaying(false);
    setProgress(0);
    // Capture final duration in case it wasn't resolved earlier
    const audio = audioRef.current;
    if (audio && isFinite(audio.duration) && audio.duration > 0) {
      setDuration(audio.duration);
    }
  };

  const handleSeek = (e: React.MouseEvent<HTMLDivElement>) => {
    const audio = audioRef.current;
    if (!audio || !duration) return;
    const rect = e.currentTarget.getBoundingClientRect();
    const ratio = (e.clientX - rect.left) / rect.width;
    audio.currentTime = ratio * duration;
    setProgress(audio.currentTime);
  };

  const formatDuration = (s: number) => {
    if (!s || !isFinite(s)) return "...";
    const m = Math.floor(s / 60);
    const sec = Math.floor(s % 60);
    return `${m}:${sec.toString().padStart(2, "0")}`;
  };

  const pct = duration > 0 ? (progress / duration) * 100 : 0;

  return (
    <div className="flex items-center gap-2.5 min-w-[180px] max-w-[240px]">
      <audio
        ref={audioRef}
        src={src}
        preload="auto"
        onEnded={handleEnded}
      >
        <source src={src} type={mime} />
      </audio>
      <button
        type="button"
        onClick={togglePlay}
        className={`flex-shrink-0 h-9 w-9 rounded-full flex items-center justify-center transition-colors ${
          isMine
            ? "bg-white/20 hover:bg-white/30 text-white"
            : "bg-primary/10 hover:bg-primary/20 text-primary"
        }`}
      >
        {playing ? <Pause className="h-4 w-4" /> : <Play className="h-4 w-4 ml-0.5" />}
      </button>
      <div className="flex-1 min-w-0">
        <div
          className={`h-1.5 rounded-full cursor-pointer ${
            isMine ? "bg-white/20" : "bg-foreground/10"
          }`}
          onClick={handleSeek}
        >
          <div
            className={`h-1.5 rounded-full transition-all ${
              isMine ? "bg-white/70" : "bg-primary/60"
            }`}
            style={{ width: `${pct}%` }}
          />
        </div>
        <span className={`text-[10px] mt-0.5 block ${
          isMine ? "text-white/50" : "text-muted-foreground"
        }`}>
          {formatDuration(playing ? progress : duration)}
        </span>
      </div>
    </div>
  );
}

// Sender avatar circle
function SenderAvatar({ name }: { name: string }) {
  const color = getAvatarColor(name);
  return (
    <div className={`flex-shrink-0 h-8 w-8 rounded-full flex items-center justify-center text-xs font-semibold ${color.bg} ${color.text}`}>
      {getInitials(name)}
    </div>
  );
}

export function ChatMessageList({ messages, currentUser, hasMore, loadingMore, onLoadMore }: ChatMessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const topSentinelRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);
  const prevScrollHeightRef = useRef<number>(0);
  const prevMessageCountRef = useRef<number>(0);
  const isLoadingMoreRef = useRef(false);

  const handleScroll = useCallback(() => {
    const el = containerRef.current;
    if (!el) return;
    const distFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    setAutoScroll(distFromBottom < 80);
  }, []);

  useEffect(() => {
    if (autoScroll && messages.length > prevMessageCountRef.current) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
    prevMessageCountRef.current = messages.length;
  }, [messages.length, autoScroll]);

  useLayoutEffect(() => {
    if (isLoadingMoreRef.current && containerRef.current) {
      const el = containerRef.current;
      const newScrollHeight = el.scrollHeight;
      const diff = newScrollHeight - prevScrollHeightRef.current;
      el.scrollTop += diff;
      isLoadingMoreRef.current = false;
    }
  }, [messages]);

  useEffect(() => {
    const sentinel = topSentinelRef.current;
    const container = containerRef.current;
    if (!sentinel || !container) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loadingMore) {
          prevScrollHeightRef.current = container.scrollHeight;
          isLoadingMoreRef.current = true;
          onLoadMore();
        }
      },
      { root: container, threshold: 0.1 }
    );

    observer.observe(sentinel);
    return () => observer.disconnect();
  }, [hasMore, loadingMore, onLoadMore]);

  // Group consecutive messages from the same sender
  const shouldShowAvatar = (msg: MessageItem, i: number): boolean => {
    if (msg.type !== "message") return false;
    if (i === 0) return true;
    const prev = messages[i - 1];
    return prev.type !== "message" || prev.senderName !== msg.senderName;
  };

  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      className="flex-1 overflow-y-auto px-3 sm:px-4 py-3 space-y-0.5"
    >
      <div ref={topSentinelRef} className="h-1" />

      {loadingMore && (
        <div className="flex justify-center py-3">
          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      )}

      {!hasMore && messages.length > 0 && (
        <div className="flex justify-center py-4 pb-6">
          <span className="text-xs text-muted-foreground/60">Beginning of conversation</span>
        </div>
      )}

      {messages.map((msg, i) => {
        if (msg.type === "join" || msg.type === "leave") {
          return (
            <div key={`sys-${i}`} className="flex justify-center py-2">
              <span className="text-[11px] text-muted-foreground/70 bg-muted/40 rounded-full px-3 py-0.5">
                <span className="font-medium">{msg.senderName}</span>{" "}
                {msg.type === "join" ? "joined" : "left"}
              </span>
            </div>
          );
        }

        const isMine = msg.senderName === currentUser;
        const media = getMediaInfo(msg.body);
        const showAvatar = shouldShowAvatar(msg, i);
        const isFirstInGroup = showAvatar;

        return (
          <div
            key={msg.id ?? `msg-${i}`}
            className={`flex items-end gap-2 ${isMine ? "flex-row-reverse" : ""} ${isFirstInGroup ? "mt-3" : "mt-0.5"}`}
          >
            {/* Avatar — only for other users, only on first message in group */}
            {!isMine ? (
              <div className="w-8 flex-shrink-0">
                {showAvatar && <SenderAvatar name={msg.senderName ?? ""} />}
              </div>
            ) : null}

            <div className={`max-w-[75%] min-w-0 ${isMine ? "items-end" : "items-start"}`}>
              {/* Sender name — only on first message in group, only for others */}
              {!isMine && showAvatar && (
                <p className="text-[11px] font-medium text-muted-foreground mb-0.5 ml-3">
                  {msg.senderName}
                </p>
              )}

              <div
                className={`rounded-2xl px-3.5 py-2 ${
                  isMine
                    ? `bg-primary text-primary-foreground ${isFirstInGroup ? "rounded-br-md" : "rounded-r-md"}`
                    : `bg-muted ${isFirstInGroup ? "rounded-bl-md" : "rounded-l-md"}`
                }`}
              >
                {media?.type === "image" ? (
                  <a href={media.src} target="_blank" rel="noopener noreferrer" className="block">
                    <img
                      src={media.src}
                      alt="shared image"
                      className="rounded-lg max-h-56 max-w-full object-contain"
                      loading="lazy"
                    />
                  </a>
                ) : media?.type === "audio" ? (
                  <AudioPlayer src={media.src} mime={media.mime} isMine={isMine} />
                ) : (
                  <p className="text-sm break-words leading-relaxed">{msg.body}</p>
                )}
                <p
                  className={`text-[10px] mt-1 ${
                    isMine ? "text-primary-foreground/50 text-right" : "text-muted-foreground/60"
                  }`}
                >
                  {formatTime(msg.createdAt)}
                </p>
              </div>
            </div>
          </div>
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
}
