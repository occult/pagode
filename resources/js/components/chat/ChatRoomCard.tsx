import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Lock, Users } from "lucide-react";
import type { ChatRoom } from "@/types/chat";

interface ChatRoomCardProps {
  room: ChatRoom;
  onClick: () => void;
}

export function ChatRoomCard({ room, onClick }: ChatRoomCardProps) {
  return (
    <Card
      className="cursor-pointer transition-all hover:shadow-md hover:border-primary/50"
      onClick={onClick}
    >
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-lg">
          <span className="truncate">{room.name}</span>
          {room.hasPassword && (
            <Lock className="h-4 w-4 text-muted-foreground flex-shrink-0" />
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="flex items-center justify-between text-sm text-muted-foreground">
          <div className="flex items-center gap-1">
            <Users className="h-3.5 w-3.5" />
            <span>{room.participantCount}</span>
          </div>
          {room.ownerName && <span>by {room.ownerName}</span>}
        </div>
      </CardContent>
    </Card>
  );
}
