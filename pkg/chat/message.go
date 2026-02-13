package chat

import "time"

// MessageType represents the type of WebSocket message.
type MessageType string

const (
	TypeMessage      MessageType = "message"
	TypeJoin         MessageType = "join"
	TypeLeave        MessageType = "leave"
	TypeTyping       MessageType = "typing"
	TypeError        MessageType = "error"
	TypeParticipants MessageType = "participants"
	TypeHistoryEnd   MessageType = "history_end"
)

// IncomingMessage is a message received from a WebSocket client.
type IncomingMessage struct {
	Type MessageType `json:"type"`
	Body string      `json:"body"`
}

// OutgoingMessage is a message sent to a WebSocket client.
type OutgoingMessage struct {
	Type         MessageType       `json:"type"`
	ID           int               `json:"id,omitempty"`
	SenderName   string            `json:"senderName,omitempty"`
	Body         string            `json:"body,omitempty"`
	CreatedAt    time.Time         `json:"createdAt,omitempty"`
	Participants []ParticipantInfo `json:"participants,omitempty"`
	HasMore      bool              `json:"hasMore,omitempty"`
}

// ParticipantInfo describes a participant in a chat room.
type ParticipantInfo struct {
	Name    string `json:"name"`
	IsOwner bool   `json:"isOwner"`
}
