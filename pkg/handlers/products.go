package handlers

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/occult/pagode/pkg/form"
	"github.com/occult/pagode/pkg/middleware"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"
	inertia "github.com/romsar/gonertia/v2"
)

type Products struct {
	Inertia *inertia.Inertia
	Payment *services.PaymentClient
	Auth    *services.AuthClient
}

func init() {
	Register(new(Products))
}

func (h *Products) Init(c *services.Container) error {
	h.Inertia = c.Inertia
	h.Payment = c.Payment
	h.Auth = c.Auth
	return nil
}

func (h *Products) Routes(g *echo.Group) {
	authGroup := g.Group("")
	authGroup.Use(middleware.RequireAuthentication)
	
	authGroup.GET("/products", h.Page).Name = routenames.Products
	authGroup.POST("/products/purchase", h.Purchase).Name = routenames.ProductsPurchase
}

func (h *Products) Page(ctx echo.Context) error {
	product := map[string]interface{}{
		"id":          "prod_your_stripe_product_id_here",
		"name":        "Premium Product",
		"description": "Access to premium features and content",
		"price":       2999, // $29.99 in cents
		"currency":    "usd",
	}

	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Products",
		inertia.Props{
			"title":               "Products",
			"product":             product,
			"stripePublishableKey": h.Payment.GetConfig().Payment.Stripe.PublishableKey,
		},
	)
}

type ProductPurchaseForm struct {
	form.Submission
	ProductID       string `form:"productId" validate:"required"`
	PaymentMethodID string `form:"paymentMethodId" validate:"required"`
	Amount          int    `form:"amount" validate:"required,min=1"`
}

func (h *Products) Purchase(ctx echo.Context) error {
	var form ProductPurchaseForm
	if err := form.Submit(ctx, &form); err != nil {
		return err
	}

	user, err := h.Auth.GetAuthenticatedUser(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	// Get or create Stripe customer
	customer, err := h.Payment.CreateOrGetCustomer(ctx, user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create customer")
	}

	// Attach payment method to customer
	_, err = h.Payment.AttachPaymentMethodToCustomer(ctx, customer, form.PaymentMethodID, true)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Failed to attach payment method: %v", err))
	}

	// Create payment intent
	paymentIntent, err := h.Payment.CreateOneTimePayment(
		ctx,
		customer,
		int64(form.Amount),
		"usd",
		fmt.Sprintf("Purchase of product %s", form.ProductID),
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Payment failed: %v", err))
	}

	// If payment succeeded, you could store purchase record in database here
	// For now, we'll just redirect with success message

	product := map[string]interface{}{
		"id":          "prod_your_stripe_product_id_here",
		"name":        "Premium Product",
		"description": "Access to premium features and content",
		"price":       2999, // $29.99 in cents
		"currency":    "usd",
	}

	return h.Inertia.Render(
		ctx.Response().Writer,
		ctx.Request(),
		"Products",
		inertia.Props{
			"title":               "Products",
			"product":             product,
			"stripePublishableKey": h.Payment.GetConfig().Payment.Stripe.PublishableKey,
			"success":             true,
			"paymentIntentId":     paymentIntent.ProviderPaymentIntentID,
		},
	)
}