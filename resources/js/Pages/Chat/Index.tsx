import { useState, type ReactNode } from "react";
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
import { Label } from "@/components/ui/label";
import { ChatRoomCard } from "@/components/chat/ChatRoomCard";
import { PasswordPrompt } from "@/components/chat/PasswordPrompt";
import { type BreadcrumbItem } from "@/types";
import type { ChatRoom } from "@/types/chat";
import { SharedProps } from "@/types/global";
import { Head, router, usePage } from "@inertiajs/react";
import { MessageCircle, Plus, UserRound } from "lucide-react";

function ChatLayout({ children, isAuth }: { children: ReactNode; isAuth: boolean }) {
  if (isAuth) {
    const breadcrumbs: BreadcrumbItem[] = [
      { title: "Dashboard", href: "/dashboard" },
      { title: "Chat", href: "/chat" },
    ];
    return <AppLayout breadcrumbs={breadcrumbs}>{children}</AppLayout>;
  }
  return <PublicLayout>{children}</PublicLayout>;
}

interface ChatIndexProps {
  rooms: ChatRoom[];
}

export default function ChatIndex({ rooms = [] }: ChatIndexProps) {
  const { auth } = usePage<SharedProps>().props;
  const isAuth = !!auth.user;
  const [showCreate, setShowCreate] = useState(false);
  const [newRoomName, setNewRoomName] = useState("");
  const [newRoomPassword, setNewRoomPassword] = useState("");
  const [creating, setCreating] = useState(false);

  // Nickname for anonymous users
  const [nickname, setNickname] = useState(() =>
    typeof window !== "undefined"
      ? localStorage.getItem("chat_nickname") || ""
      : ""
  );
  const [editingNick, setEditingNick] = useState(false);
  const [nickInput, setNickInput] = useState(nickname);

  // Nickname prompt when joining without one
  const [nicknameRoom, setNicknameRoom] = useState<ChatRoom | null>(null);
  const [nicknamePromptInput, setNicknamePromptInput] = useState("");

  // Password prompt state
  const [passwordRoom, setPasswordRoom] = useState<ChatRoom | null>(null);

  const proceedToRoom = (room: ChatRoom) => {
    if (room.hasPassword) {
      setPasswordRoom(room);
      return;
    }
    router.visit(`/chat/rooms/${room.id}`);
  };

  const handleJoinRoom = (room: ChatRoom) => {
    // Authenticated users or users with a nickname go straight through
    if (auth.user || nickname) {
      proceedToRoom(room);
      return;
    }
    // Anonymous without nickname — prompt them
    setNicknamePromptInput("");
    setNicknameRoom(room);
  };

  const handleNicknameSet = () => {
    const trimmed = nicknamePromptInput.trim();
    if (trimmed) {
      localStorage.setItem("chat_nickname", trimmed);
      setNickname(trimmed);
    }
    const room = nicknameRoom;
    setNicknameRoom(null);
    if (room) proceedToRoom(room);
  };

  const handleNicknameSkip = () => {
    const room = nicknameRoom;
    setNicknameRoom(null);
    if (room) proceedToRoom(room);
  };

  const handlePasswordSubmit = (password: string) => {
    if (!passwordRoom) return;
    sessionStorage.setItem(`chat_pwd_${passwordRoom.id}`, password);
    setPasswordRoom(null);
    router.visit(`/chat/rooms/${passwordRoom.id}`);
  };

  const handleCreateRoom = (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    router.post(
      "/chat/rooms",
      {
        name: newRoomName,
        password: newRoomPassword,
        is_public: true,
      },
      {
        onSuccess: () => {
          setShowCreate(false);
          setNewRoomName("");
          setNewRoomPassword("");
        },
        onFinish: () => setCreating(false),
      }
    );
  };

  const saveNickname = () => {
    const trimmed = nickInput.trim();
    if (trimmed) {
      localStorage.setItem("chat_nickname", trimmed);
      setNickname(trimmed);
    }
    setEditingNick(false);
  };

  return (
    <ChatLayout isAuth={isAuth}>
      <Head title="Chat" />
      <div className="flex flex-1 flex-col gap-6 rounded-xl p-4 sm:p-6">
        {/* Header */}
        <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h1 className="text-xl sm:text-2xl font-bold tracking-tight">Community Chat</h1>
            <p className="text-sm text-muted-foreground">
              Join a room to start chatting
            </p>
          </div>
          <div className="flex items-center gap-3">
            {!auth.user && (
              <div className="flex items-center gap-2">
                {editingNick ? (
                  <div className="flex items-center gap-1">
                    <Input
                      value={nickInput}
                      onChange={(e) => setNickInput(e.target.value)}
                      placeholder="Nickname"
                      className="w-32 h-8"
                      onKeyDown={(e) => e.key === "Enter" && saveNickname()}
                      autoFocus
                    />
                    <Button size="sm" variant="outline" onClick={saveNickname}>
                      Save
                    </Button>
                  </div>
                ) : (
                  <Button
                    size="sm"
                    variant="outline"
                    onClick={() => {
                      setNickInput(nickname);
                      setEditingNick(true);
                    }}
                  >
                    {nickname ? `Nick: ${nickname}` : "Set Nickname"}
                  </Button>
                )}
              </div>
            )}
            {auth.user && (
              <Button onClick={() => setShowCreate(true)}>
                <Plus className="h-4 w-4 mr-2" />
                Create Room
              </Button>
            )}
          </div>
        </div>

        {/* Room List */}
        {rooms.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-20 text-muted-foreground">
            <div className="h-16 w-16 rounded-2xl bg-muted flex items-center justify-center mb-4">
              <MessageCircle className="h-8 w-8" />
            </div>
            <p className="font-medium">No chat rooms yet</p>
            <p className="text-sm mt-1">Create one to get started</p>
          </div>
        ) : (
          <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
            {rooms.map((room) => (
              <ChatRoomCard
                key={room.id}
                room={room}
                onClick={() => handleJoinRoom(room)}
              />
            ))}
          </div>
        )}

        {/* Create Room Dialog */}
        <Dialog open={showCreate} onOpenChange={setShowCreate}>
          <DialogContent className="sm:max-w-sm">
            <DialogHeader>
              <DialogTitle>Create Chat Room</DialogTitle>
            </DialogHeader>
            <form onSubmit={handleCreateRoom} className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="room-name">Room Name</Label>
                <Input
                  id="room-name"
                  value={newRoomName}
                  onChange={(e) => setNewRoomName(e.target.value)}
                  placeholder="e.g. random"
                  maxLength={50}
                  required
                  autoFocus
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="room-pw">Password (optional)</Label>
                <Input
                  id="room-pw"
                  type="password"
                  value={newRoomPassword}
                  onChange={(e) => setNewRoomPassword(e.target.value)}
                  placeholder="Leave empty for public"
                />
              </div>
              <div className="flex justify-end gap-2">
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setShowCreate(false)}
                >
                  Cancel
                </Button>
                <Button type="submit" disabled={creating || !newRoomName.trim()}>
                  {creating ? "Creating..." : "Create"}
                </Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>

        {/* Nickname Prompt */}
        <Dialog open={!!nicknameRoom} onOpenChange={(v) => !v && setNicknameRoom(null)}>
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
                handleNicknameSet();
              }}
              className="space-y-4"
            >
              <Input
                value={nicknamePromptInput}
                onChange={(e) => setNicknamePromptInput(e.target.value)}
                placeholder="Your nickname"
                maxLength={30}
                autoFocus
              />
              <div className="flex flex-col gap-2">
                <Button type="submit" disabled={!nicknamePromptInput.trim()}>
                  Join as {nicknamePromptInput.trim() || "..."}
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
        {passwordRoom && (
          <PasswordPrompt
            open={!!passwordRoom}
            onClose={() => setPasswordRoom(null)}
            onSubmit={handlePasswordSubmit}
            roomName={passwordRoom.name}
          />
        )}
      </div>
    </ChatLayout>
  );
}
