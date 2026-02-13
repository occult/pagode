package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ChatBan holds the schema definition for the ChatBan entity.
type ChatBan struct {
	ent.Schema
}

// Fields of the ChatBan.
func (ChatBan) Fields() []ent.Field {
	return []ent.Field{
		field.String("ip_hash").
			Optional().
			MaxLen(64),
		field.String("reason").
			Optional().
			MaxLen(500),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the ChatBan.
func (ChatBan) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", ChatRoom.Type).
			Ref("bans").
			Unique().
			Required(),
		edge.From("user", User.Type).
			Ref("chat_bans").
			Unique(),
		edge.From("banned_by_user", User.Type).
			Ref("chat_bans_issued").
			Unique().
			Required(),
	}
}
