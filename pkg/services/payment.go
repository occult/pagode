package services

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/config"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/paymentcustomer"
	"github.com/occult/pagode/ent/paymentintent"
	"github.com/occult/pagode/ent/paymentmethod"
	"github.com/occult/pagode/ent/subscription"
	"github.com/occult/pagode/ent/user"
)

// PaymentProvider defines the interface for payment providers (Stripe, PayPal, etc.)
type PaymentProvider interface {
	// Customer operations
	CreateCustomer(ctx context.Context, params *CreateCustomerParams) (*CustomerResult, error)
	GetCustomer(ctx context.Context, customerID string) (*CustomerResult, error)
	UpdateCustomer(ctx context.Context, customerID string, params *UpdateCustomerParams) (*CustomerResult, error)

	// Payment intent operations (one-time payments)
	CreatePaymentIntent(ctx context.Context, params *CreatePaymentIntentParams) (*PaymentIntentResult, error)
	ConfirmPaymentIntent(ctx context.Context, paymentIntentID string, paymentMethodID string) (*PaymentIntentResult, error)
	GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResult, error)
	CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResult, error)

	// Subscription operations
	CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*SubscriptionResult, error)
	GetSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResult, error)
	UpdateSubscription(ctx context.Context, subscriptionID string, params *UpdateSubscriptionParams) (*SubscriptionResult, error)
	CancelSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResult, error)

	// Payment method operations (secure - no raw card data)
	GetPaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentMethodResult, error)
	AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*PaymentMethodResult, error)
	DetachPaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentMethodResult, error)
	ListPaymentMethods(ctx context.Context, customerID string) ([]*PaymentMethodResult, error)
	SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) (*PaymentMethodResult, error)

	// Refund operations
	CreateRefund(ctx context.Context, params *CreateRefundParams) (*RefundResult, error)
	GetRefund(ctx context.Context, refundID string) (*RefundResult, error)
}

// PaymentClient wraps the payment provider and provides high-level operations
type PaymentClient struct {
	config   *config.Config
	orm      *ent.Client
	provider PaymentProvider
}

// NewPaymentClient creates a new payment client
func NewPaymentClient(cfg *config.Config, orm *ent.Client, provider PaymentProvider) *PaymentClient {
	return &PaymentClient{
		config:   cfg,
		orm:      orm,
		provider: provider,
	}
}

// GetConfig returns the configuration
func (c *PaymentClient) GetConfig() *config.Config {
	return c.config
}

// CreateCustomerParams contains parameters for creating a customer
type CreateCustomerParams struct {
	Email    string                 `json:"email"`
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateCustomerParams contains parameters for updating a customer
type UpdateCustomerParams struct {
	Email    string                 `json:"email,omitempty"`
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// CreatePaymentIntentParams contains parameters for creating a payment intent
type CreatePaymentIntentParams struct {
	Amount      int64                  `json:"amount"`
	Currency    string                 `json:"currency"`
	CustomerID  string                 `json:"customer_id"`
	Description string                 `json:"description,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// CreateSubscriptionParams contains parameters for creating a subscription
type CreateSubscriptionParams struct {
	CustomerID      string                 `json:"customer_id"`
	PriceID         string                 `json:"price_id"`
	PaymentMethodID string                 `json:"payment_method_id,omitempty"`
	TrialPeriodDays int                    `json:"trial_period_days,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// UpdateSubscriptionParams contains parameters for updating a subscription
type UpdateSubscriptionParams struct {
	PriceID         string                 `json:"price_id,omitempty"`
	PaymentMethodID string                 `json:"payment_method_id,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// AttachPaymentMethodParams contains parameters for attaching a payment method
type AttachPaymentMethodParams struct {
	PaymentMethodID string                 `json:"payment_method_id"`
	CustomerID      string                 `json:"customer_id"`
	SetAsDefault    bool                   `json:"set_as_default,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// CreateRefundParams contains parameters for creating a refund
type CreateRefundParams struct {
	PaymentIntentID string                 `json:"payment_intent_id"`
	Amount          int64                  `json:"amount,omitempty"` // If empty, refund full amount
	Reason          string                 `json:"reason,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
}

// CustomerResult represents a customer response from the provider
type CustomerResult struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Name     string                 `json:"name,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Created  time.Time              `json:"created"`
}

// PaymentIntentResult represents a payment intent response from the provider
type PaymentIntentResult struct {
	ID           string                 `json:"id"`
	Status       string                 `json:"status"`
	Amount       int64                  `json:"amount"`
	Currency     string                 `json:"currency"`
	CustomerID   string                 `json:"customer_id"`
	Description  string                 `json:"description,omitempty"`
	ClientSecret string                 `json:"client_secret,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Created      time.Time              `json:"created"`
}

// SubscriptionResult represents a subscription response from the provider
type SubscriptionResult struct {
	ID                   string                 `json:"id"`
	Status               string                 `json:"status"`
	CustomerID           string                 `json:"customer_id"`
	PriceID              string                 `json:"price_id"`
	Amount               int64                  `json:"amount"`
	Currency             string                 `json:"currency"`
	Interval             string                 `json:"interval"`
	IntervalCount        int                    `json:"interval_count"`
	CurrentPeriodStart   time.Time              `json:"current_period_start"`
	CurrentPeriodEnd     time.Time              `json:"current_period_end"`
	TrialStart           *time.Time             `json:"trial_start,omitempty"`
	TrialEnd             *time.Time             `json:"trial_end,omitempty"`
	CanceledAt           *time.Time             `json:"canceled_at,omitempty"`
	EndedAt              *time.Time             `json:"ended_at,omitempty"`
	Metadata             map[string]interface{} `json:"metadata,omitempty"`
	Created              time.Time              `json:"created"`
}

// PaymentMethodResult represents a payment method response from the provider
type PaymentMethodResult struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	CustomerID string                 `json:"customer_id,omitempty"`
	LastFour   string                 `json:"last_four,omitempty"`
	Brand      string                 `json:"brand,omitempty"`
	ExpMonth   int                    `json:"exp_month,omitempty"`
	ExpYear    int                    `json:"exp_year,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Created    time.Time              `json:"created"`
}

// RefundResult represents a refund response from the provider
type RefundResult struct {
	ID              string                 `json:"id"`
	PaymentIntentID string                 `json:"payment_intent_id"`
	Amount          int64                  `json:"amount"`
	Currency        string                 `json:"currency"`
	Status          string                 `json:"status"`
	Reason          string                 `json:"reason,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	Created         time.Time              `json:"created"`
}

// High-level methods for the PaymentClient

// CreateOrGetCustomer creates or retrieves an existing payment customer for a user
func (c *PaymentClient) CreateOrGetCustomer(ctx echo.Context, u *ent.User) (*ent.PaymentCustomer, error) {
	// Check if customer already exists
	existingCustomer, err := c.orm.PaymentCustomer.Query().
		Where(paymentcustomer.HasUserWith(user.ID(u.ID))).
		Only(ctx.Request().Context())

	if err == nil {
		return existingCustomer, nil
	}

	// Create new customer with provider
	providerCustomer, err := c.provider.CreateCustomer(ctx.Request().Context(), &CreateCustomerParams{
		Email: u.Email,
		Name:  u.Name,
		Metadata: map[string]interface{}{
			"user_id": u.ID,
		},
	})
	if err != nil {
		return nil, err
	}

	// Save customer to database
	customer, err := c.orm.PaymentCustomer.Create().
		SetProviderCustomerID(providerCustomer.ID).
		SetProvider("stripe").
		SetEmail(providerCustomer.Email).
		SetName(providerCustomer.Name).
		SetMetadata(providerCustomer.Metadata).
		SetUser(u).
		Save(ctx.Request().Context())

	return customer, err
}

// CreateOneTimePayment creates a payment intent for a one-time payment
func (c *PaymentClient) CreateOneTimePayment(ctx echo.Context, customer *ent.PaymentCustomer, amount int64, currency, description string) (*ent.PaymentIntent, error) {
	// Create payment intent with provider
	providerPaymentIntent, err := c.provider.CreatePaymentIntent(ctx.Request().Context(), &CreatePaymentIntentParams{
		Amount:      amount,
		Currency:    currency,
		CustomerID:  customer.ProviderCustomerID,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	// Save payment intent to database
	paymentIntent, err := c.orm.PaymentIntent.Create().
		SetProviderPaymentIntentID(providerPaymentIntent.ID).
		SetProvider("stripe").
		SetStatus(paymentintent.Status(providerPaymentIntent.Status)).
		SetAmount(providerPaymentIntent.Amount).
		SetCurrency(providerPaymentIntent.Currency).
		SetDescription(providerPaymentIntent.Description).
		SetClientSecret(providerPaymentIntent.ClientSecret).
		SetMetadata(providerPaymentIntent.Metadata).
		SetCustomer(customer).
		Save(ctx.Request().Context())

	return paymentIntent, err
}

// CreateSubscription creates a new subscription
func (c *PaymentClient) CreateSubscription(ctx echo.Context, customer *ent.PaymentCustomer, priceID string, params *CreateSubscriptionParams) (*ent.Subscription, error) {
	// Create subscription with provider
	providerSubscription, err := c.provider.CreateSubscription(ctx.Request().Context(), &CreateSubscriptionParams{
		CustomerID:      customer.ProviderCustomerID,
		PriceID:         priceID,
		PaymentMethodID: params.PaymentMethodID,
		TrialPeriodDays: params.TrialPeriodDays,
		Metadata:        params.Metadata,
	})
	if err != nil {
		return nil, err
	}

	// Save subscription to database
	subscriptionBuilder := c.orm.Subscription.Create().
		SetProviderSubscriptionID(providerSubscription.ID).
		SetProvider("stripe").
		SetStatus(subscription.Status(providerSubscription.Status)).
		SetPriceID(providerSubscription.PriceID).
		SetAmount(providerSubscription.Amount).
		SetCurrency(providerSubscription.Currency).
		SetInterval(subscription.Interval(providerSubscription.Interval)).
		SetIntervalCount(providerSubscription.IntervalCount).
		SetCurrentPeriodStart(providerSubscription.CurrentPeriodStart).
		SetCurrentPeriodEnd(providerSubscription.CurrentPeriodEnd).
		SetMetadata(providerSubscription.Metadata).
		SetCustomer(customer)

	if providerSubscription.TrialStart != nil {
		subscriptionBuilder.SetTrialStart(*providerSubscription.TrialStart)
	}
	if providerSubscription.TrialEnd != nil {
		subscriptionBuilder.SetTrialEnd(*providerSubscription.TrialEnd)
	}

	subscription, err := subscriptionBuilder.Save(ctx.Request().Context())
	return subscription, err
}

// RefundPayment creates a refund for a payment intent
func (c *PaymentClient) RefundPayment(ctx echo.Context, paymentIntent *ent.PaymentIntent, amount int64, reason string) (*RefundResult, error) {
	return c.provider.CreateRefund(ctx.Request().Context(), &CreateRefundParams{
		PaymentIntentID: paymentIntent.ProviderPaymentIntentID,
		Amount:          amount,
		Reason:          reason,
	})
}

// AttachPaymentMethodToCustomer securely attaches a PaymentMethod (created by Stripe Elements) to a customer
func (c *PaymentClient) AttachPaymentMethodToCustomer(ctx echo.Context, customer *ent.PaymentCustomer, paymentMethodID string, setAsDefault bool) (*ent.PaymentMethod, error) {
	// Attach payment method to customer in Stripe
	providerPaymentMethod, err := c.provider.AttachPaymentMethod(ctx.Request().Context(), paymentMethodID, customer.ProviderCustomerID)
	if err != nil {
		return nil, err
	}

	// Set as default if requested
	if setAsDefault {
		_, err = c.provider.SetDefaultPaymentMethod(ctx.Request().Context(), customer.ProviderCustomerID, paymentMethodID)
		if err != nil {
			return nil, err
		}
	}

	// Save payment method to database (only display data)
	paymentMethod, err := c.orm.PaymentMethod.Create().
		SetProviderPaymentMethodID(providerPaymentMethod.ID).
		SetProvider("stripe").
		SetType(paymentmethod.Type(providerPaymentMethod.Type)).
		SetLastFour(providerPaymentMethod.LastFour).
		SetBrand(providerPaymentMethod.Brand).
		SetExpMonth(providerPaymentMethod.ExpMonth).
		SetExpYear(providerPaymentMethod.ExpYear).
		SetIsDefault(setAsDefault).
		SetMetadata(providerPaymentMethod.Metadata).
		SetCustomer(customer).
		Save(ctx.Request().Context())

	return paymentMethod, err
}

// GetCustomerPaymentMethods retrieves all payment methods for a customer
func (c *PaymentClient) GetCustomerPaymentMethods(ctx echo.Context, customer *ent.PaymentCustomer) ([]*ent.PaymentMethod, error) {
	return c.orm.PaymentMethod.Query().
		Where(paymentmethod.HasCustomerWith(paymentcustomer.ID(customer.ID))).
		All(ctx.Request().Context())
}

// SetDefaultPaymentMethod sets a payment method as the default for a customer
func (c *PaymentClient) SetDefaultPaymentMethod(ctx echo.Context, customer *ent.PaymentCustomer, paymentMethodID string) error {
	// First, unset all other default payment methods
	_, err := c.orm.PaymentMethod.Update().
		Where(paymentmethod.HasCustomerWith(paymentcustomer.ID(customer.ID))).
		SetIsDefault(false).
		Save(ctx.Request().Context())
	if err != nil {
		return err
	}

	// Set the new default payment method
	_, err = c.orm.PaymentMethod.Update().
		Where(paymentmethod.ProviderPaymentMethodID(paymentMethodID)).
		SetIsDefault(true).
		Save(ctx.Request().Context())
	if err != nil {
		return err
	}

	// Update in Stripe
	_, err = c.provider.SetDefaultPaymentMethod(ctx.Request().Context(), customer.ProviderCustomerID, paymentMethodID)
	return err
}

// GetCustomerSubscriptions retrieves all subscriptions for a customer
func (c *PaymentClient) GetCustomerSubscriptions(ctx echo.Context, customer *ent.PaymentCustomer) ([]*ent.Subscription, error) {
	return c.orm.Subscription.Query().
		Where(subscription.HasCustomerWith(paymentcustomer.ID(customer.ID))).
		All(ctx.Request().Context())
}

// CancelSubscription cancels a subscription
func (c *PaymentClient) CancelSubscription(ctx echo.Context, customer *ent.PaymentCustomer, subscriptionID string) error {
	// Get the subscription from database
	sub, err := c.orm.Subscription.Query().
		Where(
			subscription.ProviderSubscriptionID(subscriptionID),
			subscription.HasCustomerWith(paymentcustomer.ID(customer.ID)),
		).
		Only(ctx.Request().Context())
	if err != nil {
		return err
	}

	// Cancel with provider
	_, err = c.provider.CancelSubscription(ctx.Request().Context(), subscriptionID)
	if err != nil {
		return err
	}

	// Update subscription status in database
	_, err = c.orm.Subscription.UpdateOne(sub).
		SetStatus(subscription.StatusCanceled).
		SetCanceledAt(time.Now()).
		Save(ctx.Request().Context())
	
	return err
}

// GetCustomerPaymentIntents retrieves all payment intents for a customer
func (c *PaymentClient) GetCustomerPaymentIntents(ctx echo.Context, customer *ent.PaymentCustomer) ([]*ent.PaymentIntent, error) {
	return c.orm.PaymentIntent.Query().
		Where(paymentintent.HasCustomerWith(paymentcustomer.ID(customer.ID))).
		Order(ent.Desc(paymentintent.FieldCreatedAt)).
		All(ctx.Request().Context())
}

// ConfirmPaymentIntent confirms a payment intent with a payment method
func (c *PaymentClient) ConfirmPaymentIntent(ctx echo.Context, paymentIntent *ent.PaymentIntent, paymentMethodID string) error {
	// Confirm payment intent with provider
	providerPaymentIntent, err := c.provider.ConfirmPaymentIntent(ctx.Request().Context(), paymentIntent.ProviderPaymentIntentID, paymentMethodID)
	if err != nil {
		return err
	}

	// Update payment intent status in database
	_, err = c.orm.PaymentIntent.UpdateOne(paymentIntent).
		SetStatus(paymentintent.Status(providerPaymentIntent.Status)).
		Save(ctx.Request().Context())
	
	return err
}