package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ChatRoom holds the schema definition for the ChatRoom entity.
type ChatRoom struct {
	ent.Schema
}

// Fields of the ChatRoom.
func (ChatRoom) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Unique().
			MaxLen(50),
		field.Bool("is_public").
			Default(true),
		field.String("password_hash").
			Optional().
			Sensitive(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the ChatRoom.
func (ChatRoom) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("owned_chat_rooms").
			Unique(),
		edge.To("messages", ChatMessage.Type),
		edge.To("bans", ChatBan.Type),
	}
}
