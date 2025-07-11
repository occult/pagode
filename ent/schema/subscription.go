package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Subscription holds the schema definition for the Subscription entity.
type Subscription struct {
	ent.Schema
}

// Fields of the Subscription.
func (Subscription) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_subscription_id").
			NotEmpty().
			Unique().
			Comment("External payment provider subscription ID"),
		field.String("provider").
			NotEmpty().
			Default("stripe").
			Comment("Payment provider name"),
		field.Enum("status").
			Values("incomplete", "incomplete_expired", "trialing", "active", 
				   "past_due", "canceled", "unpaid", "paused").
			Default("incomplete").
			Comment("Subscription status from provider"),
		field.String("price_id").
			NotEmpty().
			Comment("Price/plan ID from payment provider"),
		field.Int64("amount").
			Min(0).
			Comment("Subscription amount in smallest currency unit"),
		field.String("currency").
			NotEmpty().
			Default("usd").
			Comment("Three-letter ISO currency code"),
		field.Enum("interval").
			Values("day", "week", "month", "year").
			Comment("Billing interval"),
		field.Int("interval_count").
			Default(1).
			Min(1).
			Comment("Number of intervals between billings"),
		field.Time("current_period_start").
			Optional().
			Comment("Start of current billing period"),
		field.Time("current_period_end").
			Optional().
			Comment("End of current billing period"),
		field.Time("trial_start").
			Optional().
			Comment("Trial period start"),
		field.Time("trial_end").
			Optional().
			Comment("Trial period end"),
		field.Time("canceled_at").
			Optional().
			Comment("When subscription was canceled"),
		field.Time("ended_at").
			Optional().
			Comment("When subscription ended"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional subscription data"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the Subscription.
func (Subscription) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", PaymentCustomer.Type).
			Ref("subscriptions").
			Unique().
			Required().
			Comment("Payment customer who owns this subscription"),
	}
}
