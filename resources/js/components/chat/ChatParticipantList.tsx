import type { ChatParticipant } from "@/types/chat";
import { Users } from "lucide-react";

interface ChatParticipantListProps {
  participants: ChatParticipant[];
}

export function ChatParticipantList({ participants }: ChatParticipantListProps) {
  return (
    <div className="border-l p-4 w-56 flex-shrink-0 hidden md:block">
      <div className="flex items-center gap-2 mb-3">
        <Users className="h-4 w-4 text-muted-foreground" />
        <span className="text-sm font-medium">
          Online ({participants.length})
        </span>
      </div>
      <ul className="space-y-1">
        {participants.map((p) => (
          <li key={p.name} className="flex items-center gap-2 text-sm">
            <span className="h-2 w-2 rounded-full bg-green-500 flex-shrink-0" />
            <span className="truncate">{p.name}</span>
            {p.isOwner && (
              <span className="text-xs text-muted-foreground flex-shrink-0">
                owner
              </span>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}
