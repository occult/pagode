import { type ReactNode } from "react";
import AppLayout from "@/Layouts/AppLayout";
import PublicLayout from "@/Layouts/PublicLayout";
import { ChatInput } from "@/components/chat/ChatInput";
import { ChatMessageList } from "@/components/chat/ChatMessageList";
import { ChatParticipantList } from "@/components/chat/ChatParticipantList";
import { useChat } from "@/hooks/useChat";
import { type BreadcrumbItem } from "@/types";
import type { ChatRoomDetail } from "@/types/chat";
import { SharedProps } from "@/types/global";
import { Head, usePage } from "@inertiajs/react";
import { Wifi, WifiOff } from "lucide-react";

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

  // Resolve display name
  const displayName =
    auth.user?.name ||
    (typeof window !== "undefined"
      ? localStorage.getItem("chat_nickname") || ""
      : "");

  // Resolve password from sessionStorage
  const storedPassword =
    typeof window !== "undefined"
      ? sessionStorage.getItem(`chat_pwd_${room.id}`) || undefined
      : undefined;

  const {
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
  } = useChat({
    roomId: room.id,
    name: displayName,
    password: storedPassword,
  });

  return (
    <ChatLayout isAuth={isAuth} roomName={room.name} roomId={room.id}>
      <Head title={`Chat - ${room.name}`} />
      <div className="flex flex-col rounded-xl border overflow-hidden h-[calc(100dvh-4rem)] sm:h-[calc(100dvh-8rem)]">
        {/* Top bar */}
        <div className="flex items-center justify-between border-b px-4 py-2">
          <div className="flex items-center gap-2">
            <h2 className="font-semibold">{room.name}</h2>
            {room.ownerName && (
              <span className="text-xs text-muted-foreground">
                by {room.ownerName}
              </span>
            )}
          </div>
          <div className="flex items-center gap-2">
            {connected ? (
              <Wifi className="h-4 w-4 text-green-500" />
            ) : (
              <WifiOff className="h-4 w-4 text-red-500" />
            )}
            <span className="text-xs text-muted-foreground">
              {connected ? "Connected" : "Disconnected"}
            </span>
          </div>
        </div>

        {/* Error banner */}
        {error && (
          <div className="bg-destructive/10 text-destructive text-sm px-4 py-2 border-b">
            {error}
          </div>
        )}

        {/* Main content */}
        <div className="flex flex-1 min-h-0">
          {/* Messages area */}
          <div className="flex flex-col flex-1 min-w-0">
            <ChatMessageList
              messages={messages}
              currentUser={displayName}
              hasMore={hasMore}
              loadingMore={loadingMore}
              onLoadMore={loadMore}
            />

            {/* Typing indicator */}
            {typingUsers.length > 0 && (
              <div className="px-4 pb-1 text-xs text-muted-foreground">
                {typingUsers.length === 1
                  ? `${typingUsers[0]} is typing...`
                  : `${typingUsers.join(", ")} are typing...`}
              </div>
            )}

            <ChatInput
              onSend={sendMessage}
              onTyping={sendTyping}
              disabled={!connected}
            />
          </div>

          {/* Participant sidebar */}
          <ChatParticipantList participants={participants} />
        </div>
      </div>
    </ChatLayout>
  );
}
