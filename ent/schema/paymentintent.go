package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PaymentIntent holds the schema definition for the PaymentIntent entity.
type PaymentIntent struct {
	ent.Schema
}

// Fields of the PaymentIntent.
func (PaymentIntent) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_payment_intent_id").
			NotEmpty().
			Unique().
			Comment("External payment provider payment intent ID"),
		field.String("provider").
			NotEmpty().
			Default("stripe").
			Comment("Payment provider name"),
		field.Enum("status").
			Values("requires_payment_method", "requires_confirmation", "requires_action", 
				   "processing", "requires_capture", "canceled", "succeeded").
			Default("requires_payment_method").
			Comment("Payment intent status from provider"),
		field.Int64("amount").
			Min(0).
			Comment("Amount in smallest currency unit (e.g., cents)"),
		field.String("currency").
			NotEmpty().
			Default("usd").
			Comment("Three-letter ISO currency code"),
		field.String("description").
			Optional().
			Comment("Description of the payment"),
		field.String("client_secret").
			Optional().
			Sensitive().
			Comment("Client secret for frontend payment processing"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional payment data"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the PaymentIntent.
func (PaymentIntent) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", PaymentCustomer.Type).
			Ref("payment_intents").
			Unique().
			Required().
			Comment("Payment customer who owns this payment intent"),
	}
}
