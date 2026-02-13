# WebSocket Community Chat — Design Document

**Date:** 2026-02-13
**Approach:** Hub-per-Room with gorilla/websocket

## Purpose

Add a real-time community chat to Pagode where authenticated users and anonymous visitors can communicate. Features public and private (password-protected) rooms with room-owner moderation.

## Data Model (Ent Schemas)

### ChatRoom

| Field | Type | Notes |
|-------|------|-------|
| `name` | string | Unique display name |
| `is_public` | bool | Default `true` |
| `password_hash` | string (optional) | bcrypt hash if password-protected |
| `owner_id` | UUID (optional) | FK to User. `nil` for default public room |
| `created_at` | time | Auto |

Edges: `owner` -> User (optional), `messages` -> ChatMessage (O2M), `bans` -> ChatBan (O2M)

### ChatMessage

| Field | Type | Notes |
|-------|------|-------|
| `body` | string | Max 2000 chars |
| `sender_name` | string | Display name at time of send |
| `sender_id` | UUID (optional) | FK to User if authenticated |
| `room_id` | UUID | FK to ChatRoom |
| `created_at` | time | Auto |

Edges: `room` -> ChatRoom, `sender` -> User (optional)

### ChatBan

| Field | Type | Notes |
|-------|------|-------|
| `room_id` | UUID | FK to ChatRoom |
| `banned_by` | UUID | FK to User (room owner or admin) |
| `user_id` | UUID (optional) | FK to User if banning authenticated user |
| `ip_hash` | string (optional) | sha256(ip + salt) for anonymous bans |
| `reason` | string (optional) | Ban reason |
| `created_at` | time | Auto |

Edges: `room` -> ChatRoom, `user` -> User (optional)

### ChatParticipant (in-memory only)

```go
type ChatParticipant struct {
    Conn    *websocket.Conn
    Name    string
    UserID  *uuid.UUID
    IsOwner bool
    IsAdmin bool
}
```

## WebSocket Architecture

### Hub-per-Room Pattern

**RoomManager** (singleton in service container):
- `map[uuid.UUID]*Hub` of active rooms, protected by `sync.RWMutex`
- `GetOrCreateHub(roomID)` finds or creates a Hub
- `RemoveHub(roomID)` called when a Hub empties

**Hub** (one goroutine per active room):
- `select` loop over `register`, `unregister`, `broadcast` channels
- Maintains `clients map[*ChatParticipant]bool`
- Persists messages to DB before broadcasting
- Shuts down when empty (except default public room)

**Client pumps** (per connection):
- `readPump()` reads from WebSocket, sends to Hub broadcast channel
- `writePump()` reads from per-client send channel, writes to WebSocket
- Ping/pong keepalive with 60s timeout

### Message Protocol (JSON)

Client -> Server:
```json
{"type": "message", "body": "hello everyone"}
{"type": "typing"}
```

Server -> Client:
```json
{"type": "message", "id": "uuid", "sender_name": "Felipe", "body": "hello", "created_at": "..."}
{"type": "join", "sender_name": "Felipe"}
{"type": "leave", "sender_name": "Felipe"}
{"type": "typing", "sender_name": "Felipe"}
{"type": "error", "body": "you are banned from this room"}
{"type": "participants", "participants": [{"name": "Felipe", "is_owner": false}]}
```

### Connection Flow

1. Client opens WebSocket to `/ws/chat/:roomID?name=Guest123&password=optional`
2. Handler upgrades connection
3. Extract identity: session cookie (authenticated) or `name` param (anonymous)
4. Check ChatBan — reject if banned
5. Check room password — reject if wrong
6. Register with Hub, broadcast `join`, send `participants` list
7. Send last 50 messages as history
8. readPump/writePump goroutines run until disconnect

## Routes

| Method | Path | Auth | Purpose |
|--------|------|------|---------|
| GET | `/chat` | No | Room listing (Inertia) |
| POST | `/chat/rooms` | Yes | Create room |
| GET | `/chat/rooms/:id` | No | Chat room page (Inertia) |
| GET | `/ws/chat/:id` | No | WebSocket upgrade |
| POST | `/chat/rooms/:id/ban` | Yes | Ban user |
| POST | `/chat/rooms/:id/unban` | Yes | Unban user |
| DELETE | `/chat/rooms/:id` | Yes | Delete room |

Room creation requires authentication. WebSocket is open but checks bans/passwords server-side.

## Frontend

### Pages

- `Pages/Chat/Index.tsx` — Room listing with create form and nickname input
- `Pages/Chat/Room.tsx` — Chat room with messages, input, participant sidebar

### Components

- `ChatMessageList` — Scrollable message list with auto-scroll
- `ChatInput` — Text input + send, debounced typing events
- `ChatParticipantList` — Sidebar with kick/ban for owners
- `ChatRoomCard` — Room card (name, count, lock icon)
- `PasswordPrompt` — Modal for entering room password

### Hook

`useChat(roomId, name, password?)` — manages WebSocket lifecycle, message state, actions, and auto-reconnect with exponential backoff (1s, 2s, 4s, max 30s).

## Moderation & Security

### Permissions

- **Room owners**: kick, ban, unban users in their room; delete their room
- **Admins**: same powers in any room
- Check: `isOwner || isAdmin`

### Anonymous Banning

Bans use `sha256(ip + appEncryptionKey)` hash. Not bulletproof but reasonable baseline.

### Input Validation

- Message body: max 2000 chars, trimmed, reject empty
- Room name: max 50 chars, alphanumeric + spaces + hyphens, unique
- Nickname: max 30 chars, alphanumeric + spaces
- XSS: plain text only, React default escaping

### Rate Limiting

- 10 messages per 10 seconds per connection
- Exceeding mutes for 5 seconds (no disconnect)

### Connection Limits

- 100 concurrent connections per room
- 5 concurrent connections per IP across all rooms
- 60s ping/pong timeout for stale cleanup

## Configuration

```yaml
chat:
  enabled: true
  defaultRoom: "general"
  maxMessageLength: 2000
  maxRoomsPerUser: 5
  historySize: 50
  maxConnectionsPerRoom: 100
  maxConnectionsPerIP: 5
  rateLimitMessages: 10
  rateLimitWindowSeconds: 10
```

## Startup Behavior

1. RoomManager checks for default "general" room — creates if missing
2. Spins up Hub for default room
3. Other Hubs created on-demand
4. On shutdown, all Hubs broadcast "server shutting down" and close connections

## File Organization

```
pkg/
  handlers/chat.go
  chat/
    hub.go
    manager.go
    participant.go
    message.go
ent/schema/
  chatroom.go
  chatmessage.go
  chatban.go
resources/js/
  Pages/Chat/
    Index.tsx
    Room.tsx
  components/chat/
    ChatMessageList.tsx
    ChatInput.tsx
    ChatParticipantList.tsx
    ChatRoomCard.tsx
    PasswordPrompt.tsx
  hooks/
    useChat.ts
pkg/routenames/
  ChatRooms, ChatRoomCreate, ChatRoom, ChatWebSocket, ChatBanUser, ChatUnbanUser, ChatDeleteRoom
```
