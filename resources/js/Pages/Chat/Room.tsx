import { type ReactNode, useEffect, useRef, useState } from "react";
import AppLayout from "@/Layouts/AppLayout";
import PublicLayout from "@/Layouts/PublicLayout";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { ChatInput } from "@/components/chat/ChatInput";
import { ChatMessageList } from "@/components/chat/ChatMessageList";
import { ChatParticipantList } from "@/components/chat/ChatParticipantList";
import { useChat } from "@/hooks/useChat";
import { type BreadcrumbItem } from "@/types";
import type { ChatRoomDetail } from "@/types/chat";
import { SharedProps } from "@/types/global";
import { Head, usePage } from "@inertiajs/react";
import { getAvatarColor } from "@/components/chat/avatar-colors";
import { Lock, UserRound } from "lucide-react";

function ChatLayout({ children, isAuth, roomName, roomId }: { children: ReactNode; isAuth: boolean; roomName: string; roomId: number }) {
  if (isAuth) {
    const breadcrumbs: BreadcrumbItem[] = [
      { title: "Dashboard", href: "/dashboard" },
      { title: "Chat", href: "/chat" },
      { title: roomName, href: `/chat/rooms/${roomId}` },
    ];
    return <AppLayout breadcrumbs={breadcrumbs}>{children}</AppLayout>;
  }
  return <PublicLayout>{children}</PublicLayout>;
}

interface ChatRoomProps {
  room: ChatRoomDetail;
}

export default function ChatRoom({ room }: ChatRoomProps) {
  const { auth } = usePage<SharedProps>().props;
  const isAuth = !!auth.user;

  // Resolve initial nickname
  const initialNickname =
    auth.user?.name ||
    (typeof window !== "undefined"
      ? localStorage.getItem("chat_nickname") || ""
      : "");

  // Resolve initial password from sessionStorage
  const initialPassword =
    typeof window !== "undefined"
      ? sessionStorage.getItem(`chat_pwd_${room.id}`) || ""
      : "";

  // The actual values used for the WS connection (updated when prompts are submitted)
  const [chatName, setChatName] = useState(initialNickname);
  const [chatPassword, setChatPassword] = useState(initialPassword);

  // Room owner and admins skip the password check on the backend
  const isRoomOwner = isAuth && auth.user?.id === room.ownerId;

  // Prompt states
  const [showNicknamePrompt, setShowNicknamePrompt] = useState(!isAuth && !initialNickname);
  const [showPasswordPrompt, setShowPasswordPrompt] = useState(
    room.hasPassword && !initialPassword && !isRoomOwner
  );

  // Nickname prompt input
  const [nickInput, setNickInput] = useState("");

  // Password prompt input
  const [pwdInput, setPwdInput] = useState("");

  // Only connect WS when all prompts are resolved
  const ready = !showNicknamePrompt && !showPasswordPrompt;

  const {
    messages,
    participants,
    connected,
    error,
    resetError,
    sendMessage,
    sendTyping,
    typingUsers,
    hasMore,
    loadingMore,
    loadMore,
  } = useChat({
    roomId: room.id,
    name: chatName,
    password: chatPassword || undefined,
    enabled: ready,
  });

  // If WS keeps failing on a password room, re-show the password prompt
  useEffect(() => {
    if (error && room.hasPassword && !connected && ready) {
      // Clear the bad password and ask again
      sessionStorage.removeItem(`chat_pwd_${room.id}`);
      setChatPassword("");
      setShowPasswordPrompt(true);
    }
  }, [error, room.hasPassword, room.id, connected, ready]);

  const handleNicknameSubmit = () => {
    const trimmed = nickInput.trim();
    if (trimmed) {
      localStorage.setItem("chat_nickname", trimmed);
      setChatName(trimmed);
    }
    setShowNicknamePrompt(false);
  };

  const handleNicknameSkip = () => {
    setShowNicknamePrompt(false);
  };

  const handlePasswordSubmit = () => {
    const pwd = pwdInput.trim();
    if (!pwd) return;
    sessionStorage.setItem(`chat_pwd_${room.id}`, pwd);
    resetError(); // Clear stale error before reconnecting — prevents error recovery re-triggering
    setChatPassword(pwd);
    setPwdInput("");
    setShowPasswordPrompt(false);
  };

  const roomColor = getAvatarColor(room.name);

  // Lock container height on mount so mobile browser toolbar changes don't resize it
  const containerRef = useRef<HTMLDivElement>(null);
  const [containerHeight, setContainerHeight] = useState<number | null>(null);

  useEffect(() => {
    const el = containerRef.current;
    if (!el) return;
    const top = el.getBoundingClientRect().top;
    setContainerHeight(window.innerHeight - top);
  }, []);

  return (
    <ChatLayout isAuth={isAuth} roomName={room.name} roomId={room.id}>
      <Head title={`Chat - ${room.name}`} />
      <div
        ref={containerRef}
        className="flex flex-col rounded-xl border overflow-hidden"
        style={containerHeight ? { height: containerHeight } : { height: "calc(100dvh - 8rem)" }}
      >
        {/* Top bar */}
        <div className="flex items-center justify-between border-b px-4 py-2.5 bg-card">
          <div className="flex items-center gap-3">
            <div className={`h-9 w-9 rounded-lg flex items-center justify-center text-sm font-bold ${roomColor.bg} ${roomColor.text}`}>
              {room.name.charAt(0).toUpperCase()}
            </div>
            <div>
              <h2 className="font-semibold text-sm leading-tight">{room.name}</h2>
              <p className="text-[11px] text-muted-foreground leading-tight">
                {participants.length} online
                {room.ownerName && ` \u00B7 by ${room.ownerName}`}
              </p>
            </div>
          </div>
          <div className="flex items-center gap-1.5">
            <span className={`h-2 w-2 rounded-full transition-colors ${connected ? "bg-emerald-500" : "bg-red-400 animate-pulse"}`} />
            <span className="text-[11px] text-muted-foreground">
              {connected ? "Connected" : ready ? "Reconnecting..." : "Waiting..."}
            </span>
          </div>
        </div>

        {/* Error banner */}
        {error && !showPasswordPrompt && (
          <div className="bg-destructive/10 text-destructive text-xs px-4 py-2 border-b">
            {error}
          </div>
        )}

        {/* Main content */}
        <div className="flex flex-1 min-h-0">
          {/* Messages area */}
          <div className="flex flex-col flex-1 min-w-0">
            <ChatMessageList
              messages={messages}
              currentUser={chatName}
              hasMore={hasMore}
              loadingMore={loadingMore}
              onLoadMore={loadMore}
            />

            {/* Typing indicator */}
            {typingUsers.length > 0 && (
              <div className="px-4 pb-1 text-xs text-muted-foreground flex items-center gap-1.5">
                <span className="flex gap-0.5">
                  <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/40 animate-bounce [animation-delay:0ms]" />
                  <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/40 animate-bounce [animation-delay:150ms]" />
                  <span className="h-1.5 w-1.5 rounded-full bg-muted-foreground/40 animate-bounce [animation-delay:300ms]" />
                </span>
                {typingUsers.length === 1
                  ? `${typingUsers[0]} is typing`
                  : `${typingUsers.join(", ")} are typing`}
              </div>
            )}

            <ChatInput
              onSend={sendMessage}
              onTyping={sendTyping}
              disabled={!connected}
            />
          </div>

          <ChatParticipantList participants={participants} />
        </div>
      </div>

      {/* Nickname Prompt */}
      <Dialog open={showNicknamePrompt} onOpenChange={(v) => !v && handleNicknameSkip()}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader className="items-center text-center">
            <div className="mx-auto mb-2 h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
              <UserRound className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle>Choose a nickname</DialogTitle>
            <DialogDescription>
              Pick a name so others know who you are. You can always change it later.
            </DialogDescription>
          </DialogHeader>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handleNicknameSubmit();
            }}
            className="space-y-4"
          >
            <Input
              value={nickInput}
              onChange={(e) => setNickInput(e.target.value)}
              placeholder="Your nickname"
              maxLength={30}
              autoFocus
            />
            <div className="flex flex-col gap-2">
              <Button type="submit" disabled={!nickInput.trim()}>
                Join as {nickInput.trim() || "..."}
              </Button>
              <Button
                type="button"
                variant="ghost"
                className="text-muted-foreground"
                onClick={handleNicknameSkip}
              >
                Continue without nickname
              </Button>
            </div>
          </form>
        </DialogContent>
      </Dialog>

      {/* Password Prompt */}
      <Dialog open={showPasswordPrompt && !showNicknamePrompt} onOpenChange={(v) => !v && setShowPasswordPrompt(false)}>
        <DialogContent className="sm:max-w-sm">
          <DialogHeader className="items-center text-center">
            <div className="mx-auto mb-2 h-12 w-12 rounded-full bg-primary/10 flex items-center justify-center">
              <Lock className="h-6 w-6 text-primary" />
            </div>
            <DialogTitle>Room is password-protected</DialogTitle>
            <DialogDescription>
              Enter the password to join "{room.name}".
            </DialogDescription>
          </DialogHeader>
          <form
            onSubmit={(e) => {
              e.preventDefault();
              handlePasswordSubmit();
            }}
            className="space-y-4"
          >
            <Input
              type="password"
              value={pwdInput}
              onChange={(e) => setPwdInput(e.target.value)}
              placeholder="Room password"
              autoFocus
            />
            <Button type="submit" className="w-full" disabled={!pwdInput.trim()}>
              Join Room
            </Button>
          </form>
        </DialogContent>
      </Dialog>
    </ChatLayout>
  );
}
