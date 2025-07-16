package handlers

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/ent/subscription"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Plans struct {
	Inertia *inertia.Inertia
	Payment *services.PaymentClient
	Auth    *services.AuthClient
}

func init() {
	Register(new(Plans))
}

func (h *Plans) Init(c *services.Container) error {
	h.Inertia = c.Inertia
	h.Payment = c.Payment
	h.Auth = c.Auth
	return nil
}

func (h *Plans) Routes(g *echo.Group) {
	authGroup := g.Group("")
	authGroup.Use(middleware.RequireAuthentication)
	
	authGroup.GET("/plans", h.Page).Name = routenames.Plans
	authGroup.POST("/plans/subscribe", h.Subscribe).Name = routenames.PlansSubscribe
}

func (h *Plans) Page(ctx echo.Context) error {
	// Get current user
	user, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return err
	}

	// Check if user already has a subscription
	paymentCustomer, err := h.Payment.CreateOrGetCustomer(ctx, user)
	var hasActiveSubscription bool
	if err == nil {
		// Check for active subscriptions
		subscriptions, err := h.Payment.GetCustomerSubscriptions(ctx, paymentCustomer)
		if err == nil {
			for _, sub := range subscriptions {
				if sub.Status == subscription.StatusActive {
					hasActiveSubscription = true
					break
				}
			}
		}
	}

	err = h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Plans",
		inertia.Props{
			"title":                "Choose Your Plan",
			"hasActiveSubscription": hasActiveSubscription,
			"form":                 form.Get[SubscribeForm](ctx),
			"stripePublishableKey": h.Payment.GetConfig().Payment.Stripe.PublishableKey,
			"plans": []map[string]interface{}{
				{
					"id":          "premium",
					"name":        "Premium Plan",
					"description": "Unlock all premium features",
					"price":       2900, // $29.00 in cents
					"currency":    "usd",
					"interval":    "month",
					"features": []string{
						"Advanced Analytics",
						"Priority Support", 
						"Custom Integrations",
						"Unlimited Projects",
						"API Access",
					},
					"priceId": "price_your_stripe_price_id_here", // Stripe Price ID
				},
			},
		},
	)
	if err != nil {
		handleServerErr(ctx.Response().Writer, err)
		return err
	}

	return nil
}

type SubscribeForm struct {
	form.Submission
	PlanId          string `form:"planId" validate:"required"`
	PaymentMethodId string `form:"paymentMethodId" validate:"required"`
}

func (h *Plans) Subscribe(ctx echo.Context) error {
	var input SubscribeForm

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

	// Create or get payment customer
	paymentCustomer, err := h.Payment.CreateOrGetCustomer(ctx, user)
	if err != nil {
		return fail(err, "Unable to create payment customer", h.Inertia, ctx)
	}

	// Attach the payment method to the customer
	_, err = h.Payment.AttachPaymentMethodToCustomer(ctx, paymentCustomer, input.PaymentMethodId, true)
	if err != nil {
		return fail(err, "Unable to attach payment method", h.Inertia, ctx)
	}

	// Create the subscription with real Stripe Price ID
	priceId := "price_your_stripe_price_id_here" // Stripe Price ID
	
	_, err = h.Payment.CreateSubscription(ctx, paymentCustomer, priceId, &services.CreateSubscriptionParams{
		PaymentMethodID: input.PaymentMethodId,
		Metadata: map[string]interface{}{
			"plan_id": input.PlanId,
			"user_id": fmt.Sprintf("%d", user.ID),
		},
	})
	if err != nil {
		return fail(err, "Unable to create subscription", h.Inertia, ctx)
	}

	// Add success message
	msg.Success(ctx, "Successfully subscribed to the premium plan!")

	// Return to plans page 
	return h.Page(ctx)
}