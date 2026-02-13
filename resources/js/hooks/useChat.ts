import { useCallback, useEffect, useRef, useState } from "react";
import type {
  ChatMessage,
  ChatParticipant,
  ServerMessage,
} from "@/types/chat";

interface UseChatOptions {
  roomId: number;
  name: string;
  password?: string;
}

type MessageItem = ChatMessage | { type: "join" | "leave"; senderName: string; createdAt: string };

interface UseChatReturn {
  messages: MessageItem[];
  participants: ChatParticipant[];
  connected: boolean;
  error: string | null;
  sendMessage: (body: string) => void;
  sendTyping: () => void;
  typingUsers: string[];
  hasMore: boolean;
  loadingMore: boolean;
  loadMore: () => void;
}

export function useChat({ roomId, name, password }: UseChatOptions): UseChatReturn {
  const [messages, setMessages] = useState<MessageItem[]>([]);
  const [participants, setParticipants] = useState<ChatParticipant[]>([]);
  const [connected, setConnected] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [typingUsers, setTypingUsers] = useState<string[]>([]);
  const [hasMore, setHasMore] = useState(false);
  const [loadingMore, setLoadingMore] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);
  const reconnectTimeoutRef = useRef<ReturnType<typeof setTimeout> | undefined>(undefined);
  const reconnectDelayRef = useRef(1000);
  const typingTimersRef = useRef<Map<string, ReturnType<typeof setTimeout>>>(new Map());
  const mountedRef = useRef(true);
  const messagesRef = useRef<MessageItem[]>([]);
  // Track whether we're in initial history load (before history_end)
  const historyPhaseRef = useRef(true);

  const updateMessages = useCallback((updater: MessageItem[] | ((prev: MessageItem[]) => MessageItem[])) => {
    setMessages((prev) => {
      const next = typeof updater === "function" ? updater(prev) : updater;
      messagesRef.current = next;
      return next;
    });
  }, []);

  const connect = useCallback(() => {
    if (!mountedRef.current) return;

    const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
    const params = new URLSearchParams();
    if (name) params.set("name", name);
    if (password) params.set("password", password);

    const url = `${protocol}//${window.location.host}/ws/chat/${roomId}?${params.toString()}`;
    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onopen = () => {
      if (!mountedRef.current) return;
      setConnected(true);
      setError(null);
      updateMessages([]);
      setHasMore(false);
      historyPhaseRef.current = true;
      reconnectDelayRef.current = 1000;
    };

    ws.onclose = () => {
      if (!mountedRef.current) return;
      setConnected(false);

      reconnectTimeoutRef.current = setTimeout(() => {
        reconnectDelayRef.current = Math.min(reconnectDelayRef.current * 2, 30000);
        connect();
      }, reconnectDelayRef.current);
    };

    ws.onerror = () => {
      if (!mountedRef.current) return;
      setError("Connection error");
    };

    ws.onmessage = (event) => {
      if (!mountedRef.current) return;

      const msg: ServerMessage = JSON.parse(event.data);

      switch (msg.type) {
        case "message":
          updateMessages((prev) => [...prev, msg]);
          break;
        case "join":
        case "leave":
          updateMessages((prev) => [...prev, msg]);
          break;
        case "typing": {
          const sender = msg.senderName;
          setTypingUsers((prev) =>
            prev.includes(sender) ? prev : [...prev, sender]
          );
          const existing = typingTimersRef.current.get(sender);
          if (existing) clearTimeout(existing);
          typingTimersRef.current.set(
            sender,
            setTimeout(() => {
              setTypingUsers((prev) => prev.filter((u) => u !== sender));
              typingTimersRef.current.delete(sender);
            }, 3000)
          );
          break;
        }
        case "error":
          setError(msg.body);
          break;
        case "participants":
          setParticipants(msg.participants);
          break;
        case "history_end":
          historyPhaseRef.current = false;
          setHasMore(msg.hasMore);
          break;
      }
    };
  }, [roomId, name, password, updateMessages]);

  useEffect(() => {
    mountedRef.current = true;
    connect();

    return () => {
      mountedRef.current = false;
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current);
      }
      typingTimersRef.current.forEach((timer) => clearTimeout(timer));
      typingTimersRef.current.clear();
      if (wsRef.current) {
        wsRef.current.close();
      }
    };
  }, [connect]);

  const loadMore = useCallback(async () => {
    if (loadingMore || !hasMore) return;

    // Find the oldest message ID from the ref
    let oldestId: number | undefined;
    for (const msg of messagesRef.current) {
      if (msg.type === "message" && msg.id) {
        oldestId = msg.id;
        break;
      }
    }

    if (!oldestId) return;

    setLoadingMore(true);
    try {
      const res = await fetch(`/chat/rooms/${roomId}/messages?before=${oldestId}&limit=30`);
      if (!res.ok) return;

      const data = await res.json();
      const older: MessageItem[] = (data.messages ?? []).map((m: ChatMessage) => ({
        type: "message" as const,
        id: m.id,
        senderName: m.senderName,
        body: m.body,
        createdAt: m.createdAt,
      }));

      setHasMore(data.hasMore);
      updateMessages((prev) => [...older, ...prev]);
    } catch {
      // Silently fail — user can retry by scrolling up again
    } finally {
      setLoadingMore(false);
    }
  }, [roomId, loadingMore, hasMore, updateMessages]);

  const sendMessage = useCallback((body: string) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: "message", body }));
    }
  }, []);

  const sendTyping = useCallback(() => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      wsRef.current.send(JSON.stringify({ type: "typing" }));
    }
  }, []);

  return {
    messages,
    participants,
    connected,
    error,
    sendMessage,
    sendTyping,
    typingUsers,
    hasMore,
    loadingMore,
    loadMore,
  };
}
