import { Lock, Users, ChevronRight } from "lucide-react";
import type { ChatRoom } from "@/types/chat";
import { getAvatarColor } from "./avatar-colors";

interface ChatRoomCardProps {
  room: ChatRoom;
  onClick: () => void;
}

export function ChatRoomCard({ room, onClick }: ChatRoomCardProps) {
  const color = getAvatarColor(room.name);
  const initial = room.name.charAt(0).toUpperCase();

  return (
    <button
      type="button"
      onClick={onClick}
      className="group w-full text-left rounded-xl border bg-card p-4 transition-all hover:shadow-md hover:border-primary/30 hover:-translate-y-0.5 active:translate-y-0"
    >
      <div className="flex items-center gap-3">
        {/* Room avatar */}
        <div className={`flex-shrink-0 h-11 w-11 rounded-xl flex items-center justify-center text-lg font-bold ${color.bg} ${color.text}`}>
          {initial}
        </div>

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-1.5">
            <span className="font-semibold truncate">{room.name}</span>
            {room.hasPassword && (
              <Lock className="h-3.5 w-3.5 text-muted-foreground flex-shrink-0" />
            )}
          </div>
          <div className="flex items-center gap-3 mt-0.5 text-xs text-muted-foreground">
            <span className="flex items-center gap-1">
              <Users className="h-3 w-3" />
              {room.participantCount} online
            </span>
            {room.ownerName && (
              <span className="truncate">by {room.ownerName}</span>
            )}
          </div>
        </div>

        <ChevronRight className="h-4 w-4 text-muted-foreground/40 group-hover:text-muted-foreground transition-colors flex-shrink-0" />
      </div>
    </button>
  );
}
