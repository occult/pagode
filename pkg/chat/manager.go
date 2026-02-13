package chat

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/occult/pagode/config"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/chatroom"
)

// RoomManager manages all chat room hubs.
type RoomManager struct {
	mu            sync.RWMutex
	hubs          map[int]*Hub
	orm           *ent.Client
	config        *HubConfig
	ipConns       map[string]int
	ipMu          sync.Mutex
	maxConnsPerIP int
	defaultRoomID int
	anonCounter   atomic.Int64
}

// NewRoomManager creates a new RoomManager.
func NewRoomManager(orm *ent.Client, cfg *config.ChatConfig) *RoomManager {
	return &RoomManager{
		hubs: make(map[int]*Hub),
		orm:  orm,
		config: &HubConfig{
			MaxMessageLength:       cfg.MaxMessageLength,
			HistorySize:            cfg.HistorySize,
			MaxConnectionsPerRoom:  cfg.MaxConnectionsPerRoom,
			RateLimitMessages:      cfg.RateLimitMessages,
			RateLimitWindowSeconds: cfg.RateLimitWindowSeconds,
			DefaultRoom:            cfg.DefaultRoom,
		},
		ipConns:       make(map[string]int),
		maxConnsPerIP: cfg.MaxConnectionsPerIP,
	}
}

// Init ensures the default room exists and starts its hub.
func (rm *RoomManager) Init(ctx context.Context) error {
	room, err := rm.orm.ChatRoom.Query().
		Where(chatroom.NameEQ(rm.config.DefaultRoom)).
		Only(ctx)

	if ent.IsNotFound(err) {
		room, err = rm.orm.ChatRoom.Create().
			SetName(rm.config.DefaultRoom).
			SetIsPublic(true).
			Save(ctx)
		if err != nil {
			return err
		}
		slog.Info("created default chat room", "name", rm.config.DefaultRoom)
	} else if err != nil {
		return err
	}

	rm.defaultRoomID = room.ID

	hub := NewHub(room.ID, rm.orm, rm.config, rm)
	rm.hubs[room.ID] = hub
	go hub.Run()

	return nil
}

// GetOrCreateHub returns the hub for a room, creating it if needed.
func (rm *RoomManager) GetOrCreateHub(roomID int) *Hub {
	rm.mu.RLock()
	hub, ok := rm.hubs[roomID]
	rm.mu.RUnlock()
	if ok {
		return hub
	}

	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Double-check after acquiring write lock
	if hub, ok := rm.hubs[roomID]; ok {
		return hub
	}

	hub = NewHub(roomID, rm.orm, rm.config, rm)
	rm.hubs[roomID] = hub
	go hub.Run()
	return hub
}

// RemoveHub removes a hub from the manager.
func (rm *RoomManager) RemoveHub(roomID int) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	delete(rm.hubs, roomID)
}

// TrackIPConnect increments the IP connection count. Returns false if limit reached.
func (rm *RoomManager) TrackIPConnect(ip string) bool {
	rm.ipMu.Lock()
	defer rm.ipMu.Unlock()

	if rm.ipConns[ip] >= rm.maxConnsPerIP {
		return false
	}
	rm.ipConns[ip]++
	return true
}

// TrackIPDisconnect decrements the IP connection count.
func (rm *RoomManager) TrackIPDisconnect(ip string) {
	rm.ipMu.Lock()
	defer rm.ipMu.Unlock()

	rm.ipConns[ip]--
	if rm.ipConns[ip] <= 0 {
		delete(rm.ipConns, ip)
	}
}

// GetHubClientCount returns the number of connected clients for a room.
func (rm *RoomManager) GetHubClientCount(roomID int) int {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	if hub, ok := rm.hubs[roomID]; ok {
		return hub.ClientCount()
	}
	return 0
}

// NextAnonName returns a unique anonymous name like "Anonymous #1".
func (rm *RoomManager) NextAnonName() string {
	n := rm.anonCounter.Add(1)
	return fmt.Sprintf("Anonymous #%d", n)
}

// Shutdown shuts down all hubs.
func (rm *RoomManager) Shutdown() {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for _, hub := range rm.hubs {
		hub.Shutdown()
	}
	rm.hubs = make(map[int]*Hub)
}
