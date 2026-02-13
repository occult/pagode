export interface ChatRoom {
  id: number;
  name: string;
  isPublic: boolean;
  hasPassword: boolean;
  ownerName: string;
  participantCount: number;
}

export interface ChatRoomDetail {
  id: number;
  name: string;
  isPublic: boolean;
  hasPassword: boolean;
  ownerName: string;
  ownerId: number;
}

export interface ChatMessage {
  type: "message";
  id: number;
  senderName: string;
  body: string;
  createdAt: string;
}

export interface ChatJoinLeave {
  type: "join" | "leave";
  senderName: string;
  createdAt: string;
}

export interface ChatTyping {
  type: "typing";
  senderName: string;
}

export interface ChatError {
  type: "error";
  body: string;
}

export interface ChatParticipantList {
  type: "participants";
  participants: ChatParticipant[];
}

export interface ChatHistoryEnd {
  type: "history_end";
  hasMore: boolean;
}

export interface ChatParticipant {
  name: string;
  isOwner: boolean;
}

export type ServerMessage =
  | ChatMessage
  | ChatJoinLeave
  | ChatTyping
  | ChatError
  | ChatParticipantList
  | ChatHistoryEnd;
