package chat

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/chatmessage"
	"github.com/occult/pagode/ent/chatroom"
)

// HubConfig holds configuration values needed by a Hub.
type HubConfig struct {
	MaxMessageLength       int
	HistorySize            int
	MaxConnectionsPerRoom  int
	RateLimitMessages      int
	RateLimitWindowSeconds int
	DefaultRoom            string
}

// BroadcastMsg wraps a message from a participant.
type BroadcastMsg struct {
	Sender  *Participant
	Message *IncomingMessage
}

// Hub maintains the set of active participants for a single chat room.
type Hub struct {
	RoomID     int
	clients    map[*Participant]bool
	broadcast  chan *BroadcastMsg
	register   chan *Participant
	unregister chan *Participant
	orm        *ent.Client
	config     *HubConfig
	manager    *RoomManager
	done       chan struct{}
}

// NewHub creates a new Hub for a room.
func NewHub(roomID int, orm *ent.Client, cfg *HubConfig, mgr *RoomManager) *Hub {
	return &Hub{
		RoomID:     roomID,
		clients:    make(map[*Participant]bool),
		broadcast:  make(chan *BroadcastMsg, 256),
		register:   make(chan *Participant),
		unregister: make(chan *Participant),
		orm:        orm,
		config:     cfg,
		manager:    mgr,
		done:       make(chan struct{}),
	}
}

// Register adds a participant to the hub.
func (h *Hub) Register(p *Participant) {
	h.register <- p
}

// Run starts the hub's event loop in a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case participant := <-h.register:
			h.clients[participant] = true

			// Broadcast join
			h.broadcastJSON(OutgoingMessage{
				Type:       TypeJoin,
				SenderName: participant.Name,
				CreatedAt:  time.Now(),
			})

			// Send history to new participant
			h.sendHistory(participant)

			// Send updated participant list to all
			h.sendParticipants()

		case participant := <-h.unregister:
			if _, ok := h.clients[participant]; ok {
				delete(h.clients, participant)
				close(participant.Send)

				// Track IP disconnect
				h.manager.TrackIPDisconnect(participant.IP)

				// Broadcast leave
				h.broadcastJSON(OutgoingMessage{
					Type:       TypeLeave,
					SenderName: participant.Name,
					CreatedAt:  time.Now(),
				})

				// Send updated participant list
				h.sendParticipants()

				// Self-destruct if empty and not default room
				if len(h.clients) == 0 {
					room, err := h.orm.ChatRoom.Get(context.Background(), h.RoomID)
					if err == nil && room.Name != h.config.DefaultRoom {
						h.manager.RemoveHub(h.RoomID)
						return
					}
				}
			}

		case bcast := <-h.broadcast:
			switch bcast.Message.Type {
			case TypeMessage:
				// Persist to DB
				msgBuilder := h.orm.ChatMessage.Create().
					SetBody(bcast.Message.Body).
					SetSenderName(bcast.Sender.Name).
					SetRoomID(h.RoomID)

				if bcast.Sender.UserID > 0 {
					msgBuilder.SetSenderID(bcast.Sender.UserID)
				}

				saved, err := msgBuilder.Save(context.Background())
				if err != nil {
					slog.Error("failed to persist chat message", "err", err)
				}

				out := OutgoingMessage{
					Type:       TypeMessage,
					SenderName: bcast.Sender.Name,
					Body:       bcast.Message.Body,
					CreatedAt:  time.Now(),
				}
				if saved != nil {
					out.ID = saved.ID
					out.CreatedAt = saved.CreatedAt
				}

				data, _ := json.Marshal(out)
				for client := range h.clients {
					select {
					case client.Send <- data:
					default:
						close(client.Send)
						delete(h.clients, client)
					}
				}

			case TypeTyping:
				out := OutgoingMessage{
					Type:       TypeTyping,
					SenderName: bcast.Sender.Name,
				}
				data, _ := json.Marshal(out)
				for client := range h.clients {
					if client != bcast.Sender {
						select {
						case client.Send <- data:
						default:
						}
					}
				}
			}

		case <-h.done:
			for client := range h.clients {
				close(client.Send)
				delete(h.clients, client)
			}
			return
		}
	}
}

// broadcastJSON marshals and sends a message to all clients.
func (h *Hub) broadcastJSON(msg OutgoingMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(h.clients, client)
		}
	}
}

// sendHistory sends the last N messages to a participant, followed by a history_end marker.
func (h *Hub) sendHistory(p *Participant) {
	// Fetch one extra to detect if there are more
	messages, err := h.orm.ChatMessage.Query().
		Where(chatmessage.HasRoomWith(chatroom.IDEQ(h.RoomID))).
		Order(ent.Desc(chatmessage.FieldID)).
		Limit(h.config.HistorySize + 1).
		All(context.Background())
	if err != nil {
		slog.Error("failed to load chat history", "err", err)
		return
	}

	hasMore := len(messages) > h.config.HistorySize
	if hasMore {
		messages = messages[:h.config.HistorySize]
	}

	// Reverse to send oldest first
	for i := len(messages) - 1; i >= 0; i-- {
		m := messages[i]
		data, _ := json.Marshal(OutgoingMessage{
			Type:       TypeMessage,
			ID:         m.ID,
			SenderName: m.SenderName,
			Body:       m.Body,
			CreatedAt:  m.CreatedAt,
		})
		select {
		case p.Send <- data:
		default:
			return
		}
	}

	// Send history_end so frontend knows initial load is done
	end, _ := json.Marshal(OutgoingMessage{
		Type:    TypeHistoryEnd,
		HasMore: hasMore,
	})
	select {
	case p.Send <- end:
	default:
	}
}

// sendParticipants sends the current participant list to all clients.
func (h *Hub) sendParticipants() {
	participants := make([]ParticipantInfo, 0, len(h.clients))
	for client := range h.clients {
		participants = append(participants, ParticipantInfo{
			Name:    client.Name,
			IsOwner: client.IsOwner,
		})
	}

	data, _ := json.Marshal(OutgoingMessage{
		Type:         TypeParticipants,
		Participants: participants,
	})

	for client := range h.clients {
		select {
		case client.Send <- data:
		default:
		}
	}
}

// Shutdown gracefully shuts down the hub.
func (h *Hub) Shutdown() {
	select {
	case h.done <- struct{}{}:
	default:
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	return len(h.clients)
}

// KickUser removes a user from the hub by user ID.
func (h *Hub) KickUser(userID int) {
	for client := range h.clients {
		if client.UserID == userID {
			h.unregister <- client
			return
		}
	}
}
