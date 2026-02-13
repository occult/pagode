package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ChatMessage holds the schema definition for the ChatMessage entity.
type ChatMessage struct {
	ent.Schema
}

// Fields of the ChatMessage.
func (ChatMessage) Fields() []ent.Field {
	return []ent.Field{
		field.String("body").
			NotEmpty().
			MaxLen(2000),
		field.String("sender_name").
			NotEmpty().
			MaxLen(30),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the ChatMessage.
func (ChatMessage) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("room", ChatRoom.Type).
			Ref("messages").
			Unique().
			Required(),
		edge.From("sender", User.Type).
			Ref("chat_messages").
			Unique(),
	}
}
