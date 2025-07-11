package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PaymentCustomer holds the schema definition for the PaymentCustomer entity.
type PaymentCustomer struct {
	ent.Schema
}

// Fields of the PaymentCustomer.
func (PaymentCustomer) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_customer_id").
			NotEmpty().
			Comment("External payment provider customer ID (e.g., Stripe customer ID)"),
		field.String("provider").
			NotEmpty().
			Default("stripe").
			Comment("Payment provider name"),
		field.String("email").
			NotEmpty().
			Comment("Customer email from payment provider"),
		field.String("name").
			Optional().
			Comment("Customer name from payment provider"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional customer data from payment provider"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the PaymentCustomer.
func (PaymentCustomer) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Comment("User associated with this payment customer"),
		edge.To("payment_intents", PaymentIntent.Type).
			Comment("Payment intents for this customer"),
		edge.To("subscriptions", Subscription.Type).
			Comment("Subscriptions for this customer"),
		edge.To("payment_methods", PaymentMethod.Type).
			Comment("Payment methods for this customer"),
	}
}
