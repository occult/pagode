package handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Billing struct {
	Inertia *inertia.Inertia
	Payment *services.PaymentClient
	Auth    *services.AuthClient
}

func init() {
	Register(new(Billing))
}

func (h *Billing) Init(c *services.Container) error {
	h.Inertia = c.Inertia
	h.Payment = c.Payment
	h.Auth = c.Auth
	return nil
}

func (h *Billing) Routes(g *echo.Group) {
	authGroup := g.Group("")
	authGroup.Use(middleware.RequireAuthentication)
	
	authGroup.GET("/billing", h.Page).Name = routenames.Billing
	authGroup.POST("/billing/cancel", h.CancelSubscription).Name = routenames.BillingCancel
}

func (h *Billing) Page(ctx echo.Context) error {
	// Get current user
	user, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	// Get payment customer
	paymentCustomer, err := h.Payment.CreateOrGetCustomer(ctx, user)
	if err != nil {
		return err
	}

	// Get user's subscriptions
	subscriptions, err := h.Payment.GetCustomerSubscriptions(ctx, paymentCustomer)
	if err != nil {
		// If no subscriptions found, continue with empty slice
		subscriptions = []*ent.Subscription{}
	}

	// Get payment methods
	paymentMethods, err := h.Payment.GetCustomerPaymentMethods(ctx, paymentCustomer)
	if err != nil {
		// If no payment methods found, continue with empty slice
		paymentMethods = []*ent.PaymentMethod{}
	}

	// Convert subscriptions to the format expected by React
	subscriptionData := make([]map[string]interface{}, len(subscriptions))
	for i, sub := range subscriptions {
		subscriptionData[i] = map[string]interface{}{
			"id":                  sub.ProviderSubscriptionID,
			"status":              string(sub.Status),
			"currentPeriodStart":  sub.CurrentPeriodStart.Format("2006-01-02T15:04:05Z07:00"),
			"currentPeriodEnd":    sub.CurrentPeriodEnd.Format("2006-01-02T15:04:05Z07:00"),
			"amount":              sub.Amount,
			"currency":            sub.Currency,
			"interval":            string(sub.Interval),
			"metadata":            sub.Metadata,
		}
	}

	// Convert payment methods to the format expected by React
	paymentMethodData := make([]map[string]interface{}, len(paymentMethods))
	for i, pm := range paymentMethods {
		paymentMethodData[i] = map[string]interface{}{
			"id":          pm.ProviderPaymentMethodID,
			"type":        string(pm.Type),
			"brand":       pm.Brand,
			"lastFour":    pm.LastFour,
			"expiryMonth": pm.ExpMonth,
			"expiryYear":  pm.ExpYear,
			"isDefault":   pm.IsDefault,
		}
	}

	err = h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Billing",
		inertia.Props{
			"title":          "Billing & Subscription",
			"subscriptions":  subscriptionData,
			"paymentMethods": paymentMethodData,
			"user":           user,
			"form":           form.Get[CancelSubscriptionForm](ctx),
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

type CancelSubscriptionForm struct {
	form.Submission
	SubscriptionID string `form:"subscriptionId" validate:"required"`
}

func (h *Billing) CancelSubscription(ctx echo.Context) error {
	var input CancelSubscriptionForm

	err := form.Submit(ctx, &input)

	switch err.(type) {
	case nil:
	case validator.ValidationErrors:
		return h.Page(ctx)
	default:
		return err
	}

	// Get current user
	user, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return fail(err, "Unable to get authenticated user", h.Inertia, ctx)
	}

	// Get payment customer
	paymentCustomer, err := h.Payment.CreateOrGetCustomer(ctx, user)
	if err != nil {
		return fail(err, "Unable to get payment customer", h.Inertia, ctx)
	}

	// Cancel the subscription
	err = h.Payment.CancelSubscription(ctx, paymentCustomer, input.SubscriptionID)
	if err != nil {
		return fail(err, "Unable to cancel subscription", h.Inertia, ctx)
	}

	// Add success message
	msg.Success(ctx, "Subscription cancelled successfully!")

	// Return to billing page 
	return h.Page(ctx)
}