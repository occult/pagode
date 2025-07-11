package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// PaymentMethod holds the schema definition for the PaymentMethod entity.
type PaymentMethod struct {
	ent.Schema
}

// Fields of the PaymentMethod.
func (PaymentMethod) Fields() []ent.Field {
	return []ent.Field{
		field.String("provider_payment_method_id").
			NotEmpty().
			Unique().
			Comment("External payment provider payment method ID"),
		field.String("provider").
			NotEmpty().
			Default("stripe").
			Comment("Payment provider name"),
		field.Enum("type").
			Values("card", "bank_account", "wallet").
			Default("card").
			Comment("Payment method type"),
		field.String("last_four").
			Optional().
			Comment("Last four digits of card/account"),
		field.String("brand").
			Optional().
			Comment("Card brand (visa, mastercard, etc.)"),
		field.Int("exp_month").
			Optional().
			Min(1).
			Max(12).
			Comment("Card expiration month"),
		field.Int("exp_year").
			Optional().
			Comment("Card expiration year"),
		field.Bool("is_default").
			Default(false).
			Comment("Whether this is the default payment method"),
		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Additional payment method data"),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the PaymentMethod.
func (PaymentMethod) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("customer", PaymentCustomer.Type).
			Ref("payment_methods").
			Unique().
			Required().
			Comment("Payment customer who owns this payment method"),
	}
}
