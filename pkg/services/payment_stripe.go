package services

import (
	"context"
	"time"

	"github.com/occult/pagode/config"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/customer"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/paymentmethod"
	"github.com/stripe/stripe-go/v82/refund"
	"github.com/stripe/stripe-go/v82/subscription"
)

// StripeProvider implements the PaymentProvider interface for Stripe
type StripeProvider struct {
	config *config.Config
}

// NewStripeProvider creates a new Stripe payment provider
func NewStripeProvider(cfg *config.Config) *StripeProvider {
	stripe.Key = cfg.Payment.Stripe.SecretKey
	return &StripeProvider{
		config: cfg,
	}
}

// CreateCustomer creates a new customer in Stripe
func (s *StripeProvider) CreateCustomer(ctx context.Context, params *CreateCustomerParams) (*CustomerResult, error) {
	stripeParams := &stripe.CustomerParams{
		Email: stripe.String(params.Email),
	}

	if params.Name != "" {
		stripeParams.Name = stripe.String(params.Name)
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	cust, err := customer.New(stripeParams)
	if err != nil {
		return nil, err
	}

	return &CustomerResult{
		ID:       cust.ID,
		Email:    cust.Email,
		Name:     cust.Name,
		Metadata: convertStripeMetadata(cust.Metadata),
		Created:  time.Unix(cust.Created, 0),
	}, nil
}

// GetCustomer retrieves a customer from Stripe
func (s *StripeProvider) GetCustomer(ctx context.Context, customerID string) (*CustomerResult, error) {
	cust, err := customer.Get(customerID, nil)
	if err != nil {
		return nil, err
	}

	return &CustomerResult{
		ID:       cust.ID,
		Email:    cust.Email,
		Name:     cust.Name,
		Metadata: convertStripeMetadata(cust.Metadata),
		Created:  time.Unix(cust.Created, 0),
	}, nil
}

// UpdateCustomer updates a customer in Stripe
func (s *StripeProvider) UpdateCustomer(ctx context.Context, customerID string, params *UpdateCustomerParams) (*CustomerResult, error) {
	stripeParams := &stripe.CustomerParams{}

	if params.Email != "" {
		stripeParams.Email = stripe.String(params.Email)
	}

	if params.Name != "" {
		stripeParams.Name = stripe.String(params.Name)
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	cust, err := customer.Update(customerID, stripeParams)
	if err != nil {
		return nil, err
	}

	return &CustomerResult{
		ID:       cust.ID,
		Email:    cust.Email,
		Name:     cust.Name,
		Metadata: convertStripeMetadata(cust.Metadata),
		Created:  time.Unix(cust.Created, 0),
	}, nil
}

// CreatePaymentIntent creates a new payment intent in Stripe
func (s *StripeProvider) CreatePaymentIntent(ctx context.Context, params *CreatePaymentIntentParams) (*PaymentIntentResult, error) {
	stripeParams := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(params.Amount),
		Currency: stripe.String(params.Currency),
		Customer: stripe.String(params.CustomerID),
	}

	if params.Description != "" {
		stripeParams.Description = stripe.String(params.Description)
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	pi, err := paymentintent.New(stripeParams)
	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		ID:           pi.ID,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		CustomerID:   pi.Customer.ID,
		Description:  pi.Description,
		ClientSecret: pi.ClientSecret,
		Metadata:     convertStripeMetadata(pi.Metadata),
		Created:      time.Unix(pi.Created, 0),
	}, nil
}

// ConfirmPaymentIntent confirms a payment intent in Stripe
func (s *StripeProvider) ConfirmPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResult, error) {
	pi, err := paymentintent.Confirm(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		ID:           pi.ID,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		CustomerID:   pi.Customer.ID,
		Description:  pi.Description,
		ClientSecret: pi.ClientSecret,
		Metadata:     convertStripeMetadata(pi.Metadata),
		Created:      time.Unix(pi.Created, 0),
	}, nil
}

// GetPaymentIntent retrieves a payment intent from Stripe
func (s *StripeProvider) GetPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResult, error) {
	pi, err := paymentintent.Get(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		ID:           pi.ID,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		CustomerID:   pi.Customer.ID,
		Description:  pi.Description,
		ClientSecret: pi.ClientSecret,
		Metadata:     convertStripeMetadata(pi.Metadata),
		Created:      time.Unix(pi.Created, 0),
	}, nil
}

// CancelPaymentIntent cancels a payment intent in Stripe
func (s *StripeProvider) CancelPaymentIntent(ctx context.Context, paymentIntentID string) (*PaymentIntentResult, error) {
	pi, err := paymentintent.Cancel(paymentIntentID, nil)
	if err != nil {
		return nil, err
	}

	return &PaymentIntentResult{
		ID:           pi.ID,
		Status:       string(pi.Status),
		Amount:       pi.Amount,
		Currency:     string(pi.Currency),
		CustomerID:   pi.Customer.ID,
		Description:  pi.Description,
		ClientSecret: pi.ClientSecret,
		Metadata:     convertStripeMetadata(pi.Metadata),
		Created:      time.Unix(pi.Created, 0),
	}, nil
}

// CreateSubscription creates a new sub in Stripe
func (s *StripeProvider) CreateSubscription(ctx context.Context, params *CreateSubscriptionParams) (*SubscriptionResult, error) {
	stripeParams := &stripe.SubscriptionParams{
		Customer: stripe.String(params.CustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(params.PriceID),
			},
		},
	}

	if params.PaymentMethodID != "" {
		stripeParams.DefaultPaymentMethod = stripe.String(params.PaymentMethodID)
	}

	if params.TrialPeriodDays > 0 {
		stripeParams.TrialPeriodDays = stripe.Int64(int64(params.TrialPeriodDays))
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	sub, err := subscription.New(stripeParams)
	if err != nil {
		return nil, err
	}

	result := &SubscriptionResult{
		ID:         sub.ID,
		Status:     string(sub.Status),
		CustomerID: sub.Customer.ID,
		Metadata:   convertStripeMetadata(sub.Metadata),
		Created:    time.Unix(sub.Created, 0),
	}

	// Get pricing information and period information from the first item
	if len(sub.Items.Data) > 0 {
		item := sub.Items.Data[0]
		result.PriceID = item.Price.ID
		result.Amount = item.Price.UnitAmount
		result.Currency = string(item.Price.Currency)
		result.Interval = string(item.Price.Recurring.Interval)
		result.IntervalCount = int(item.Price.Recurring.IntervalCount)

		// Period information is now at the subscription item level
		result.CurrentPeriodStart = time.Unix(item.CurrentPeriodStart, 0)
		result.CurrentPeriodEnd = time.Unix(item.CurrentPeriodEnd, 0)
	}

	if sub.TrialStart != 0 {
		trialStart := time.Unix(sub.TrialStart, 0)
		result.TrialStart = &trialStart
	}

	if sub.TrialEnd != 0 {
		trialEnd := time.Unix(sub.TrialEnd, 0)
		result.TrialEnd = &trialEnd
	}

	if sub.CanceledAt != 0 {
		canceledAt := time.Unix(sub.CanceledAt, 0)
		result.CanceledAt = &canceledAt
	}

	if sub.EndedAt != 0 {
		endedAt := time.Unix(sub.EndedAt, 0)
		result.EndedAt = &endedAt
	}

	return result, nil
}

// GetSubscription retrieves a sub from Stripe
func (s *StripeProvider) GetSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResult, error) {
	sub, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		return nil, err
	}

	result := &SubscriptionResult{
		ID:         sub.ID,
		Status:     string(sub.Status),
		CustomerID: sub.Customer.ID,
		Metadata:   convertStripeMetadata(sub.Metadata),
		Created:    time.Unix(sub.Created, 0),
	}

	// Get pricing information and period information from the first item
	if len(sub.Items.Data) > 0 {
		item := sub.Items.Data[0]
		result.PriceID = item.Price.ID
		result.Amount = item.Price.UnitAmount
		result.Currency = string(item.Price.Currency)
		result.Interval = string(item.Price.Recurring.Interval)
		result.IntervalCount = int(item.Price.Recurring.IntervalCount)

		// Period information is now at the subscription item level
		result.CurrentPeriodStart = time.Unix(item.CurrentPeriodStart, 0)
		result.CurrentPeriodEnd = time.Unix(item.CurrentPeriodEnd, 0)
	}

	if sub.TrialStart != 0 {
		trialStart := time.Unix(sub.TrialStart, 0)
		result.TrialStart = &trialStart
	}

	if sub.TrialEnd != 0 {
		trialEnd := time.Unix(sub.TrialEnd, 0)
		result.TrialEnd = &trialEnd
	}

	if sub.CanceledAt != 0 {
		canceledAt := time.Unix(sub.CanceledAt, 0)
		result.CanceledAt = &canceledAt
	}

	if sub.EndedAt != 0 {
		endedAt := time.Unix(sub.EndedAt, 0)
		result.EndedAt = &endedAt
	}

	return result, nil
}

// UpdateSubscription updates a sub in Stripe
func (s *StripeProvider) UpdateSubscription(ctx context.Context, subscriptionID string, params *UpdateSubscriptionParams) (*SubscriptionResult, error) {
	stripeParams := &stripe.SubscriptionParams{}

	if params.PriceID != "" {
		// To update price, we need to get the current sub first
		currentSub, err := subscription.Get(subscriptionID, nil)
		if err != nil {
			return nil, err
		}

		if len(currentSub.Items.Data) > 0 {
			stripeParams.Items = []*stripe.SubscriptionItemsParams{
				{
					ID:    stripe.String(currentSub.Items.Data[0].ID),
					Price: stripe.String(params.PriceID),
				},
			}
		}
	}

	if params.PaymentMethodID != "" {
		stripeParams.DefaultPaymentMethod = stripe.String(params.PaymentMethodID)
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	sub, err := subscription.Update(subscriptionID, stripeParams)
	if err != nil {
		return nil, err
	}

	result := &SubscriptionResult{
		ID:         sub.ID,
		Status:     string(sub.Status),
		CustomerID: sub.Customer.ID,
		Metadata:   convertStripeMetadata(sub.Metadata),
		Created:    time.Unix(sub.Created, 0),
	}

	// Get pricing information and period information from the first item
	if len(sub.Items.Data) > 0 {
		item := sub.Items.Data[0]
		result.PriceID = item.Price.ID
		result.Amount = item.Price.UnitAmount
		result.Currency = string(item.Price.Currency)
		result.Interval = string(item.Price.Recurring.Interval)
		result.IntervalCount = int(item.Price.Recurring.IntervalCount)

		// Period information is now at the subscription item level
		result.CurrentPeriodStart = time.Unix(item.CurrentPeriodStart, 0)
		result.CurrentPeriodEnd = time.Unix(item.CurrentPeriodEnd, 0)
	}

	return result, nil
}

// CancelSubscription cancels a sub in Stripe
func (s *StripeProvider) CancelSubscription(ctx context.Context, subscriptionID string) (*SubscriptionResult, error) {
	sub, err := subscription.Cancel(subscriptionID, nil)
	if err != nil {
		return nil, err
	}

	result := &SubscriptionResult{
		ID:         sub.ID,
		Status:     string(sub.Status),
		CustomerID: sub.Customer.ID,
		Metadata:   convertStripeMetadata(sub.Metadata),
		Created:    time.Unix(sub.Created, 0),
	}

	// Get pricing information and period information from the first item
	if len(sub.Items.Data) > 0 {
		item := sub.Items.Data[0]
		result.PriceID = item.Price.ID
		result.Amount = item.Price.UnitAmount
		result.Currency = string(item.Price.Currency)
		result.Interval = string(item.Price.Recurring.Interval)
		result.IntervalCount = int(item.Price.Recurring.IntervalCount)

		// Period information is now at the subscription item level
		result.CurrentPeriodStart = time.Unix(item.CurrentPeriodStart, 0)
		result.CurrentPeriodEnd = time.Unix(item.CurrentPeriodEnd, 0)
	}

	if sub.CanceledAt != 0 {
		canceledAt := time.Unix(sub.CanceledAt, 0)
		result.CanceledAt = &canceledAt
	}

	if sub.EndedAt != 0 {
		endedAt := time.Unix(sub.EndedAt, 0)
		result.EndedAt = &endedAt
	}

	return result, nil
}

// GetPaymentMethod retrieves a payment method from Stripe
func (s *StripeProvider) GetPaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentMethodResult, error) {
	pm, err := paymentmethod.Get(paymentMethodID, nil)
	if err != nil {
		return nil, err
	}

	result := &PaymentMethodResult{
		ID:       pm.ID,
		Type:     string(pm.Type),
		Metadata: convertStripeMetadata(pm.Metadata),
		Created:  time.Unix(pm.Created, 0),
	}

	if pm.Customer != nil {
		result.CustomerID = pm.Customer.ID
	}

	if pm.Card != nil {
		result.LastFour = pm.Card.Last4
		result.Brand = string(pm.Card.Brand)
		result.ExpMonth = int(pm.Card.ExpMonth)
		result.ExpYear = int(pm.Card.ExpYear)
	}

	return result, nil
}

// AttachPaymentMethod attaches a payment method to a customer in Stripe
func (s *StripeProvider) AttachPaymentMethod(ctx context.Context, paymentMethodID, customerID string) (*PaymentMethodResult, error) {
	pm, err := paymentmethod.Attach(paymentMethodID, &stripe.PaymentMethodAttachParams{
		Customer: stripe.String(customerID),
	})
	if err != nil {
		return nil, err
	}

	result := &PaymentMethodResult{
		ID:         pm.ID,
		Type:       string(pm.Type),
		CustomerID: pm.Customer.ID,
		Metadata:   convertStripeMetadata(pm.Metadata),
		Created:    time.Unix(pm.Created, 0),
	}

	if pm.Card != nil {
		result.LastFour = pm.Card.Last4
		result.Brand = string(pm.Card.Brand)
		result.ExpMonth = int(pm.Card.ExpMonth)
		result.ExpYear = int(pm.Card.ExpYear)
	}

	return result, nil
}

// DetachPaymentMethod detaches a payment method from a customer in Stripe
func (s *StripeProvider) DetachPaymentMethod(ctx context.Context, paymentMethodID string) (*PaymentMethodResult, error) {
	pm, err := paymentmethod.Detach(paymentMethodID, nil)
	if err != nil {
		return nil, err
	}

	result := &PaymentMethodResult{
		ID:       pm.ID,
		Type:     string(pm.Type),
		Metadata: convertStripeMetadata(pm.Metadata),
		Created:  time.Unix(pm.Created, 0),
	}

	if pm.Card != nil {
		result.LastFour = pm.Card.Last4
		result.Brand = string(pm.Card.Brand)
		result.ExpMonth = int(pm.Card.ExpMonth)
		result.ExpYear = int(pm.Card.ExpYear)
	}

	return result, nil
}

// ListPaymentMethods lists payment methods for a customer in Stripe
func (s *StripeProvider) ListPaymentMethods(ctx context.Context, customerID string) ([]*PaymentMethodResult, error) {
	params := &stripe.PaymentMethodListParams{
		Customer: stripe.String(customerID),
		Type:     stripe.String("card"),
	}

	iter := paymentmethod.List(params)
	var results []*PaymentMethodResult

	for iter.Next() {
		pm := iter.PaymentMethod()
		result := &PaymentMethodResult{
			ID:         pm.ID,
			Type:       string(pm.Type),
			CustomerID: pm.Customer.ID,
			Metadata:   convertStripeMetadata(pm.Metadata),
			Created:    time.Unix(pm.Created, 0),
		}

		if pm.Card != nil {
			result.LastFour = pm.Card.Last4
			result.Brand = string(pm.Card.Brand)
			result.ExpMonth = int(pm.Card.ExpMonth)
			result.ExpYear = int(pm.Card.ExpYear)
		}

		results = append(results, result)
	}

	return results, iter.Err()
}

// SetDefaultPaymentMethod sets a payment method as the default for a customer
func (s *StripeProvider) SetDefaultPaymentMethod(ctx context.Context, customerID, paymentMethodID string) (*PaymentMethodResult, error) {
	// Update the customer's default payment method
	_, err := customer.Update(customerID, &stripe.CustomerParams{
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(paymentMethodID),
		},
	})
	if err != nil {
		return nil, err
	}

	// Return the payment method details
	return s.GetPaymentMethod(ctx, paymentMethodID)
}

// CreateRefund creates a refund in Stripe
func (s *StripeProvider) CreateRefund(ctx context.Context, params *CreateRefundParams) (*RefundResult, error) {
	stripeParams := &stripe.RefundParams{
		PaymentIntent: stripe.String(params.PaymentIntentID),
	}

	if params.Amount > 0 {
		stripeParams.Amount = stripe.Int64(params.Amount)
	}

	if params.Reason != "" {
		stripeParams.Reason = stripe.String(params.Reason)
	}

	if params.Metadata != nil {
		stripeParams.Metadata = make(map[string]string)
		for k, v := range params.Metadata {
			if str, ok := v.(string); ok {
				stripeParams.Metadata[k] = str
			}
		}
	}

	r, err := refund.New(stripeParams)
	if err != nil {
		return nil, err
	}

	return &RefundResult{
		ID:              r.ID,
		PaymentIntentID: r.PaymentIntent.ID,
		Amount:          r.Amount,
		Currency:        string(r.Currency),
		Status:          string(r.Status),
		Reason:          string(r.Reason),
		Metadata:        convertStripeMetadata(r.Metadata),
		Created:         time.Unix(r.Created, 0),
	}, nil
}

// GetRefund retrieves a refund from Stripe
func (s *StripeProvider) GetRefund(ctx context.Context, refundID string) (*RefundResult, error) {
	r, err := refund.Get(refundID, nil)
	if err != nil {
		return nil, err
	}

	return &RefundResult{
		ID:              r.ID,
		PaymentIntentID: r.PaymentIntent.ID,
		Amount:          r.Amount,
		Currency:        string(r.Currency),
		Status:          string(r.Status),
		Reason:          string(r.Reason),
		Metadata:        convertStripeMetadata(r.Metadata),
		Created:         time.Unix(r.Created, 0),
	}, nil
}

// Helper function to convert Stripe metadata to map[string]interface{}
func convertStripeMetadata(metadata map[string]string) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range metadata {
		result[k] = v
	}
	return result
}

