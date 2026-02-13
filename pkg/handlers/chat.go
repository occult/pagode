package handlers

import (
	"fmt"
	"io"
	"log/slog"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/chatban"
	"github.com/occult/pagode/ent/chatmessage"
	"github.com/occult/pagode/ent/chatroom"
	"github.com/occult/pagode/ent/user"
	"github.com/occult/pagode/pkg/chat"
	appctx "github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
	"golang.org/x/crypto/bcrypt"
)

type Chat struct {
	orm     *ent.Client
	chat    *chat.RoomManager
	auth    *services.AuthClient
	config  *config.Config
	Inertia *inertia.Inertia
}

func init() {
	Register(new(Chat))

	// Register MIME types for audio/video formats that may not be in Go's default set.
	// This ensures the static file server sends the correct Content-Type header.
	mime.AddExtensionType(".webm", "audio/webm")
	mime.AddExtensionType(".m4a", "audio/mp4")
	mime.AddExtensionType(".aac", "audio/aac")
	mime.AddExtensionType(".ogg", "audio/ogg")
}

func (h *Chat) Init(c *services.Container) error {
	h.orm = c.ORM
	h.chat = c.Chat
	h.auth = c.Auth
	h.config = c.Config
	h.Inertia = c.Inertia
	return nil
}

func (h *Chat) Routes(g *echo.Group) {
	if !h.config.Chat.Enabled {
		return
	}

	g.GET("/chat", h.Index).Name = routenames.ChatRooms
	g.GET("/chat/rooms/:id", h.Room).Name = routenames.ChatRoom
	g.GET("/chat/rooms/:id/messages", h.Messages)
	g.POST("/chat/upload", h.UploadFile)

	authGroup := g.Group("")
	authGroup.Use(middleware.RequireAuthentication)
	authGroup.POST("/chat/rooms", h.CreateRoom).Name = routenames.ChatRoomCreate
	authGroup.POST("/chat/rooms/:id/ban", h.BanUser).Name = routenames.ChatBanUser
	authGroup.POST("/chat/rooms/:id/unban", h.UnbanUser).Name = routenames.ChatUnbanUser
	authGroup.DELETE("/chat/rooms/:id", h.DeleteRoom).Name = routenames.ChatDeleteRoom
}

func (h *Chat) RoutesWS(wsG *echo.Group) {
	if !h.config.Chat.Enabled {
		return
	}
	wsG.GET("/ws/chat/:id", h.WebSocket).Name = routenames.ChatWebSocket
}

func (h *Chat) Index(ctx echo.Context) error {
	rooms, err := h.orm.ChatRoom.Query().
		Order(ent.Asc(chatroom.FieldCreatedAt)).
		WithOwner().
		All(ctx.Request().Context())
	if err != nil {
		return err
	}

	type roomProps struct {
		ID               int    `json:"id"`
		Name             string `json:"name"`
		IsPublic         bool   `json:"isPublic"`
		HasPassword      bool   `json:"hasPassword"`
		OwnerName        string `json:"ownerName"`
		ParticipantCount int    `json:"participantCount"`
	}

	roomList := make([]roomProps, 0, len(rooms))
	for _, r := range rooms {
		rp := roomProps{
			ID:               r.ID,
			Name:             r.Name,
			IsPublic:         r.IsPublic,
			HasPassword:      r.PasswordHash != "",
			ParticipantCount: h.chat.GetHubClientCount(r.ID),
		}
		if owner, err := r.Edges.OwnerOrErr(); err == nil {
			rp.OwnerName = owner.Name
		}
		roomList = append(roomList, rp)
	}

	err = h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Chat/Index",
		inertia.Props{
			"rooms": roomList,
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}
	return nil
}

func (h *Chat) Room(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	room, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(id)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	type roomDetailProps struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		IsPublic    bool   `json:"isPublic"`
		HasPassword bool   `json:"hasPassword"`
		OwnerName   string `json:"ownerName"`
		OwnerID     int    `json:"ownerId"`
	}

	props := roomDetailProps{
		ID:          room.ID,
		Name:        room.Name,
		IsPublic:    room.IsPublic,
		HasPassword: room.PasswordHash != "",
	}
	if owner, err := room.Edges.OwnerOrErr(); err == nil {
		props.OwnerName = owner.Name
		props.OwnerID = owner.ID
	}

	err = h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Chat/Room",
		inertia.Props{
			"room": props,
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}
	return nil
}

// Messages returns paginated messages for a room (JSON).
// Query params: before (message ID cursor), limit (default 30, max 100).
func (h *Chat) Messages(ctx echo.Context) error {
	roomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Verify room exists
	exists, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(roomID)).
		Exist(ctx.Request().Context())
	if err != nil || !exists {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	limit := 30
	if l := ctx.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	query := h.orm.ChatMessage.Query().
		Where(chatmessage.HasRoomWith(chatroom.IDEQ(roomID))).
		Order(ent.Desc(chatmessage.FieldID)).
		Limit(limit + 1) // Fetch one extra to detect if there are more

	if before := ctx.QueryParam("before"); before != "" {
		if beforeID, err := strconv.Atoi(before); err == nil {
			query = query.Where(chatmessage.IDLT(beforeID))
		}
	}

	messages, err := query.All(ctx.Request().Context())
	if err != nil {
		return err
	}

	hasMore := len(messages) > limit
	if hasMore {
		messages = messages[:limit]
	}

	type messageJSON struct {
		ID         int    `json:"id"`
		SenderName string `json:"senderName"`
		Body       string `json:"body"`
		CreatedAt  string `json:"createdAt"`
	}

	// Reverse so oldest first
	result := make([]messageJSON, len(messages))
	for i, m := range messages {
		result[len(messages)-1-i] = messageJSON{
			ID:         m.ID,
			SenderName: m.SenderName,
			Body:       m.Body,
			CreatedAt:  m.CreatedAt.Format(time.RFC3339Nano),
		}
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"messages": result,
		"hasMore":  hasMore,
	})
}

// CreateRoomForm represents the form for creating a chat room.
type CreateRoomForm struct {
	form.Submission
	Name     string `form:"name" validate:"required,max=50"`
	Password string `form:"password"`
	IsPublic bool   `form:"is_public"`
}

func (h *Chat) CreateRoom(ctx echo.Context) error {
	var input CreateRoomForm
	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		msg.Danger(ctx, "Please fix the errors below.")
		return h.Index(ctx)
	default:
		return err
	}

	authUser, err := h.auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	// Check room limit per user
	count, err := h.orm.ChatRoom.Query().
		Where(chatroom.HasOwnerWith(user.IDEQ(authUser.ID))).
		Count(ctx.Request().Context())
	if err != nil {
		return err
	}
	if count >= h.config.Chat.MaxRoomsPerUser {
		msg.Danger(ctx, "You have reached the maximum number of rooms.")
		return h.Index(ctx)
	}

	builder := h.orm.ChatRoom.Create().
		SetName(input.Name).
		SetIsPublic(input.IsPublic).
		SetOwnerID(authUser.ID)

	if input.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		builder.SetPasswordHash(string(hash))
	}

	_, err = builder.Save(ctx.Request().Context())
	if err != nil {
		if ent.IsConstraintError(err) {
			msg.Danger(ctx, "A room with that name already exists.")
			return h.Index(ctx)
		}
		return err
	}

	msg.Success(ctx, "Room created successfully!")
	h.Inertia.Back(ctx.Response().Writer, ctx.Request())
	return nil
}

// BanUserForm represents the form for banning a user.
type BanUserForm struct {
	form.Submission
	UserID int    `form:"user_id" validate:"required"`
	Reason string `form:"reason" validate:"max=500"`
}

func (h *Chat) BanUser(ctx echo.Context) error {
	roomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	var input BanUserForm
	if err := form.Submit(ctx, &input); err != nil {
		return err
	}

	authUser, err := h.auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	// Verify owner or admin
	room, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(roomID)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	owner, _ := room.Edges.OwnerOrErr()
	if !authUser.Admin && (owner == nil || owner.ID != authUser.ID) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	_, err = h.orm.ChatBan.Create().
		SetRoomID(roomID).
		SetUserID(input.UserID).
		SetBannedByUserID(authUser.ID).
		SetReason(input.Reason).
		Save(ctx.Request().Context())
	if err != nil {
		return err
	}

	// Kick from hub
	hub := h.chat.GetOrCreateHub(roomID)
	hub.KickUser(input.UserID)

	msg.Success(ctx, "User banned successfully.")
	h.Inertia.Back(ctx.Response().Writer, ctx.Request())
	return nil
}

func (h *Chat) UnbanUser(ctx echo.Context) error {
	roomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	userID, err := strconv.Atoi(ctx.FormValue("user_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	authUser, err := h.auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	// Verify owner or admin
	room, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(roomID)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	owner, _ := room.Edges.OwnerOrErr()
	if !authUser.Admin && (owner == nil || owner.ID != authUser.ID) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Delete ban
	_, err = h.orm.ChatBan.Delete().
		Where(
			chatban.HasRoomWith(chatroom.IDEQ(roomID)),
			chatban.HasUserWith(user.IDEQ(userID)),
		).
		Exec(ctx.Request().Context())
	if err != nil {
		return err
	}

	msg.Success(ctx, "User unbanned successfully.")
	h.Inertia.Back(ctx.Response().Writer, ctx.Request())
	return nil
}

func (h *Chat) DeleteRoom(ctx echo.Context) error {
	roomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	authUser, err := h.auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	room, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(roomID)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Don't allow deleting the default room
	if room.Name == h.config.Chat.DefaultRoom {
		msg.Danger(ctx, "Cannot delete the default room.")
		h.Inertia.Back(ctx.Response().Writer, ctx.Request())
		return nil
	}

	owner, _ := room.Edges.OwnerOrErr()
	if !authUser.Admin && (owner == nil || owner.ID != authUser.ID) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	// Shutdown hub
	hub := h.chat.GetOrCreateHub(roomID)
	hub.Shutdown()
	h.chat.RemoveHub(roomID)

	// Delete messages, bans, then room
	h.orm.ChatMessage.Delete().
		Where(chatmessage.HasRoomWith(chatroom.IDEQ(roomID))).
		Exec(ctx.Request().Context())
	h.orm.ChatBan.Delete().
		Where(chatban.HasRoomWith(chatroom.IDEQ(roomID))).
		Exec(ctx.Request().Context())
	h.orm.ChatRoom.DeleteOneID(roomID).Exec(ctx.Request().Context())

	msg.Success(ctx, "Room deleted successfully.")
	return ctx.Redirect(http.StatusSeeOther, "/chat")
}

// uploadRateLimiter tracks per-IP upload counts for rate limiting.
type uploadRateLimiter struct {
	mu      sync.Mutex
	buckets map[string]*uploadBucket
}

type uploadBucket struct {
	count  int
	expiry time.Time
}

var uploadLimiter = &uploadRateLimiter{
	buckets: make(map[string]*uploadBucket),
}

// allow checks if the IP is within the upload rate limit (10 per minute).
func (l *uploadRateLimiter) allow(ip string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	b, ok := l.buckets[ip]
	if !ok || now.After(b.expiry) {
		l.buckets[ip] = &uploadBucket{count: 1, expiry: now.Add(time.Minute)}
		return true
	}
	if b.count >= 10 {
		return false
	}
	b.count++
	return true
}

// allowedMIMETypes maps detected MIME types to file extensions.
// We validate by reading actual file bytes, not the Content-Type header.
var allowedMIMETypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/gif":       ".gif",
	"image/webp":      ".webp",
	"audio/webm":      ".webm",
	"video/webm":      ".webm",
	"audio/ogg":       ".ogg",
	"audio/mpeg":      ".mp3",
	"application/ogg": ".ogg",
	"audio/mp4":       ".m4a",
	"audio/aac":       ".aac",
	"audio/x-m4a":     ".m4a",
	"video/mp4":       ".m4a",
}

// maxUploadSize is the maximum file upload size (5MB).
const maxUploadSize = 5 * 1024 * 1024

func (h *Chat) wsUpgrader() websocket.Upgrader {
	appHost := h.config.App.Host
	return websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			if origin == "" {
				return true
			}
			// Parse the configured app host to compare
			parsed, err := url.Parse(appHost)
			if err != nil {
				return false
			}
			originParsed, err := url.Parse(origin)
			if err != nil {
				return false
			}
			// Allow if hostname matches (covers localhost and tunnel domains)
			return originParsed.Hostname() == parsed.Hostname() || originParsed.Hostname() == r.Host || originParsed.Host == r.Host
		},
	}
}

func (h *Chat) WebSocket(ctx echo.Context) error {
	roomID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Verify room exists
	room, err := h.orm.ChatRoom.Query().
		Where(chatroom.IDEQ(roomID)).
		WithOwner().
		Only(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	// Get user info
	var userID int
	var userName string
	var isOwner bool
	var isAdmin bool

	if u := ctx.Get(appctx.AuthenticatedUserKey); u != nil {
		authUser := u.(*ent.User)
		userID = authUser.ID
		userName = authUser.Name
		isAdmin = authUser.Admin
		if owner, err := room.Edges.OwnerOrErr(); err == nil {
			isOwner = owner.ID == authUser.ID
		}

		// Check ban
		banned, _ := h.orm.ChatBan.Query().
			Where(
				chatban.HasRoomWith(chatroom.IDEQ(roomID)),
				chatban.HasUserWith(user.IDEQ(authUser.ID)),
			).
			Exist(ctx.Request().Context())
		if banned {
			return echo.NewHTTPError(http.StatusForbidden, "you are banned from this room")
		}
	}

	// Get name from query for anonymous users
	if userName == "" {
		userName = ctx.QueryParam("name")
		if userName == "" {
			userName = h.chat.NextAnonName()
		}
		if len(userName) > 30 {
			userName = userName[:30]
		}
	}

	// Password check
	if room.PasswordHash != "" {
		password := ctx.QueryParam("password")
		if err := bcrypt.CompareHashAndPassword([]byte(room.PasswordHash), []byte(password)); err != nil {
			return echo.NewHTTPError(http.StatusForbidden, "incorrect room password")
		}
	}

	// IP limit check
	ip := ctx.RealIP()
	if !h.chat.TrackIPConnect(ip) {
		return echo.NewHTTPError(http.StatusTooManyRequests, "too many connections from your IP")
	}

	// Upgrade to WebSocket
	up := h.wsUpgrader()
	conn, err := up.Upgrade(ctx.Response(), ctx.Request(), nil)
	if err != nil {
		h.chat.TrackIPDisconnect(ip)
		slog.Error("websocket upgrade failed", "err", err)
		return nil
	}

	hub := h.chat.GetOrCreateHub(roomID)

	participant := &chat.Participant{
		Conn:    conn,
		Name:    userName,
		UserID:  userID,
		IsOwner: isOwner,
		IsAdmin: isAdmin,
		Send:    make(chan []byte, 256),
		Hub:     hub,
		IP:      ip,
	}
	participant.SetRateLimits(h.config.Chat.RateLimitMessages, h.config.Chat.RateLimitWindowSeconds)

	hub.Register(participant)

	go participant.ReadPump()
	go participant.WritePump()

	return nil
}

func (h *Chat) UploadFile(ctx echo.Context) error {
	// Rate limit: 10 uploads per minute per IP
	ip := ctx.RealIP()
	if !uploadLimiter.allow(ip) {
		return echo.NewHTTPError(http.StatusTooManyRequests, "upload rate limit exceeded, try again in a minute")
	}

	// Try "file" field first, fall back to "image" for backwards compatibility
	file, err := ctx.FormFile("file")
	if err != nil {
		file, err = ctx.FormFile("image")
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "no file provided")
		}
	}

	if file.Size > maxUploadSize {
		return echo.NewHTTPError(http.StatusBadRequest, "file too large (max 5MB)")
	}

	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Read first 512 bytes to detect actual content type via magic bytes
	header := make([]byte, 512)
	n, err := src.Read(header)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "could not read file")
	}
	header = header[:n]

	detectedType := http.DetectContentType(header)

	// Strip parameters like ";codecs=..." for lookup (e.g. "audio/mp4;codecs=mp4a.40.2" -> "audio/mp4")
	baseType := detectedType
	if idx := strings.IndexByte(baseType, ';'); idx != -1 {
		baseType = strings.TrimSpace(baseType[:idx])
	}

	// For WebM/Ogg/MP4 audio, http.DetectContentType often returns "application/octet-stream"
	// or "video/webm", so also check the Content-Type header for audio formats
	ext, ok := allowedMIMETypes[baseType]
	if !ok || baseType == "application/octet-stream" {
		// Fall back to the client-provided Content-Type for formats
		// that http.DetectContentType doesn't recognize well
		clientType := file.Header.Get("Content-Type")
		// Strip parameters from client type too (iOS sends "audio/mp4;codecs=mp4a.40.2")
		if idx := strings.IndexByte(clientType, ';'); idx != -1 {
			clientType = strings.TrimSpace(clientType[:idx])
		}
		if clientExt, clientOK := allowedMIMETypes[clientType]; clientOK {
			ext = clientExt
			ok = true
		}
		if !ok {
			slog.Warn("unsupported upload type", "detected", detectedType, "client", file.Header.Get("Content-Type"))
			return echo.NewHTTPError(http.StatusBadRequest, "unsupported file type")
		}
	}

	// Seek back to start so we copy the full file
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return err
	}

	uploadDir := filepath.Join("static", "chat-uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return err
	}

	// Generate safe filename
	safeName := strings.Map(func(r rune) rune {
		if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' || r == '.' {
			return r
		}
		return '_'
	}, filepath.Base(file.Filename))
	filename := fmt.Sprintf("%d-%s%s", time.Now().UnixNano(), safeName, ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return err
	}

	fileURL := fmt.Sprintf("/files/chat-uploads/%s", filename)
	return ctx.JSON(http.StatusOK, map[string]string{"url": fileURL})
}
