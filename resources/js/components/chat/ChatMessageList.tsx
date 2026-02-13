import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";
import { Loader2 } from "lucide-react";

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

export function ChatMessageList({ messages, currentUser, hasMore, loadingMore, onLoadMore }: ChatMessageListProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const bottomRef = useRef<HTMLDivElement>(null);
  const topSentinelRef = useRef<HTMLDivElement>(null);
  const [autoScroll, setAutoScroll] = useState(true);
  // Track scroll height before prepending older messages
  const prevScrollHeightRef = useRef<number>(0);
  const prevMessageCountRef = useRef<number>(0);
  const isLoadingMoreRef = useRef(false);

  // Detect when user scrolls away from bottom
  const handleScroll = useCallback(() => {
    const el = containerRef.current;
    if (!el) return;
    const distFromBottom = el.scrollHeight - el.scrollTop - el.clientHeight;
    setAutoScroll(distFromBottom < 80);
  }, []);

  // Auto-scroll to bottom for new messages (only if user is near bottom)
  useEffect(() => {
    if (autoScroll && messages.length > prevMessageCountRef.current) {
      bottomRef.current?.scrollIntoView({ behavior: "smooth" });
    }
    prevMessageCountRef.current = messages.length;
  }, [messages.length, autoScroll]);

  // Preserve scroll position when older messages are prepended
  useLayoutEffect(() => {
    if (isLoadingMoreRef.current && containerRef.current) {
      const el = containerRef.current;
      const newScrollHeight = el.scrollHeight;
      const diff = newScrollHeight - prevScrollHeightRef.current;
      el.scrollTop += diff;
      isLoadingMoreRef.current = false;
    }
  }, [messages]);

  // Intersection observer for scroll-to-top auto-loading
  useEffect(() => {
    const sentinel = topSentinelRef.current;
    const container = containerRef.current;
    if (!sentinel || !container) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting && hasMore && !loadingMore) {
          // Save scroll height before loading more
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

  return (
    <div
      ref={containerRef}
      onScroll={handleScroll}
      className="flex-1 overflow-y-auto p-4 space-y-3"
    >
      {/* Top sentinel for infinite scroll */}
      <div ref={topSentinelRef} className="h-1" />

      {/* Loading indicator */}
      {loadingMore && (
        <div className="flex justify-center py-2">
          <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
        </div>
      )}

      {/* "No more messages" indicator */}
      {!hasMore && messages.length > 0 && (
        <div className="flex justify-center py-2">
          <span className="text-xs text-muted-foreground">Beginning of conversation</span>
        </div>
      )}

      {messages.map((msg, i) => {
        if (msg.type === "join" || msg.type === "leave") {
          return (
            <div key={`sys-${i}`} className="flex justify-center py-1">
              <span className="text-xs text-muted-foreground bg-muted/50 rounded-full px-3 py-0.5">
                <span className="font-medium">{msg.senderName}</span>{" "}
                {msg.type === "join" ? "joined" : "left"}
              </span>
            </div>
          );
        }

        const isMine = msg.senderName === currentUser;
        const media = getMediaInfo(msg.body);

        return (
          <div
            key={msg.id ?? `msg-${i}`}
            className={`flex ${isMine ? "justify-end" : "justify-start"}`}
          >
            <div
              className={`max-w-[75%] rounded-2xl px-4 py-2 ${
                isMine
                  ? "bg-primary text-white rounded-br-sm"
                  : "bg-muted rounded-bl-sm"
              }`}
            >
              {!isMine && (
                <p className="text-xs font-semibold mb-0.5 opacity-70">
                  {msg.senderName}
                </p>
              )}
              {media?.type === "image" ? (
                <a href={media.src} target="_blank" rel="noopener noreferrer">
                  <img
                    src={media.src}
                    alt="shared image"
                    className="rounded-lg max-h-64 max-w-full object-contain"
                    loading="lazy"
                  />
                </a>
              ) : media?.type === "audio" ? (
                <audio controls preload="metadata" className="max-w-full min-w-[200px]">
                  <source src={media.src} type={media.mime} />
                </audio>
              ) : (
                <p className="text-sm break-words">{msg.body}</p>
              )}
              <p
                className={`text-[10px] mt-1 ${
                  isMine ? "text-white/60 text-right" : "text-muted-foreground"
                }`}
              >
                {formatTime(msg.createdAt)}
              </p>
            </div>
          </div>
        );
      })}
      <div ref={bottomRef} />
    </div>
  );
}
