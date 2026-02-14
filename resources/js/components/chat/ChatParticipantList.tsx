import type { ChatParticipant } from "@/types/chat";
import { Users, Crown } from "lucide-react";
import { getAvatarColor, getInitials } from "./avatar-colors";

interface ChatParticipantListProps {
  participants: ChatParticipant[];
}

export function ChatParticipantList({ participants }: ChatParticipantListProps) {
  return (
    <div className="border-l p-4 w-56 flex-shrink-0 hidden md:flex flex-col">
      <div className="flex items-center gap-2 mb-4">
        <Users className="h-4 w-4 text-muted-foreground" />
        <span className="text-sm font-medium">
          Online ({participants.length})
        </span>
      </div>
      <ul className="space-y-1.5 overflow-y-auto flex-1">
        {participants.map((p) => {
          const color = getAvatarColor(p.name);
          return (
            <li key={p.name} className="flex items-center gap-2.5 px-2 py-1.5 rounded-lg hover:bg-muted/50 transition-colors">
              <div className="relative flex-shrink-0">
                <div className={`h-7 w-7 rounded-full flex items-center justify-center text-[10px] font-semibold ${color.bg} ${color.text}`}>
                  {getInitials(p.name)}
                </div>
                <span className="absolute -bottom-0.5 -right-0.5 h-2.5 w-2.5 rounded-full bg-emerald-500 border-2 border-background" />
              </div>
              <span className="text-sm truncate flex-1">{p.name}</span>
              {p.isOwner && (
                <Crown className="h-3.5 w-3.5 text-amber-500 flex-shrink-0" />
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}
