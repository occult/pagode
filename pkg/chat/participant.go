package chat

import (
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 54 * time.Second
	maxMessageSize = 4096
)

// Participant represents a connected WebSocket client.
type Participant struct {
	Conn    *websocket.Conn
	Name    string
	UserID  int
	IsOwner bool
	IsAdmin bool
	Send    chan []byte
	Hub     *Hub
	IP      string

	mu            sync.Mutex
	messageCount  int
	windowStart   time.Time
	rateLimitMsgs int
	rateLimitSecs int
}

// SetRateLimits sets the rate limit parameters for a participant.
func (p *Participant) SetRateLimits(messages, windowSeconds int) {
	p.rateLimitMsgs = messages
	p.rateLimitSecs = windowSeconds
}

// ReadPump pumps messages from the WebSocket connection to the hub.
func (p *Participant) ReadPump() {
	defer func() {
		p.Hub.unregister <- p
		p.Conn.Close()
	}()

	p.Conn.SetReadLimit(maxMessageSize)
	p.Conn.SetReadDeadline(time.Now().Add(pongWait))
	p.Conn.SetPongHandler(func(string) error {
		p.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, raw, err := p.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				slog.Warn("websocket read error", "err", err)
			}
			return
		}

		var msg IncomingMessage
		if err := json.Unmarshal(raw, &msg); err != nil {
			p.sendError("invalid message format")
			continue
		}

		switch msg.Type {
		case TypeMessage:
			body := strings.TrimSpace(msg.Body)
			if body == "" {
				continue
			}
			if len(body) > p.Hub.config.MaxMessageLength {
				body = body[:p.Hub.config.MaxMessageLength]
			}
			if !p.checkRateLimit() {
				p.sendError("you are sending messages too fast")
				continue
			}
			p.Hub.broadcast <- &BroadcastMsg{
				Sender:  p,
				Message: &IncomingMessage{Type: TypeMessage, Body: body},
			}
		case TypeTyping:
			p.Hub.broadcast <- &BroadcastMsg{
				Sender:  p,
				Message: &IncomingMessage{Type: TypeTyping},
			}
		}
	}
}

// WritePump pumps messages from the hub to the WebSocket connection.
func (p *Participant) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-p.Send:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := p.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			p.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (p *Participant) checkRateLimit() bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	window := time.Duration(p.rateLimitSecs) * time.Second

	if now.Sub(p.windowStart) > window {
		p.messageCount = 0
		p.windowStart = now
	}

	p.messageCount++
	return p.messageCount <= p.rateLimitMsgs
}

func (p *Participant) sendError(msg string) {
	data, _ := json.Marshal(OutgoingMessage{
		Type: TypeError,
		Body: msg,
	})
	select {
	case p.Send <- data:
	default:
	}
}
