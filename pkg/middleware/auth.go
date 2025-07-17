package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/occult/pagode/ent"
	"github.com/occult/pagode/ent/paymentcustomer"
	"github.com/occult/pagode/ent/paymentintent"
	"github.com/occult/pagode/ent/subscription"
	entuser "github.com/occult/pagode/ent/user"
	"github.com/occult/pagode/pkg/context"
	"github.com/occult/pagode/pkg/log"
	"github.com/occult/pagode/pkg/msg"
	"github.com/occult/pagode/pkg/routenames"
	"github.com/occult/pagode/pkg/services"

	"github.com/labstack/echo/v4"
)

// LoadAuthenticatedUser loads the authenticated user, if one, and stores in context.
func LoadAuthenticatedUser(authClient *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			u, err := authClient.GetAuthenticatedUser(c)
			switch err.(type) {
			case *ent.NotFoundError:
				log.Ctx(c).Warn("auth user not found")
			case services.NotAuthenticatedError:
			case nil:
				c.Set(context.AuthenticatedUserKey, u)
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error querying for authenticated user: %v", err),
				)
			}

			return next(c)
		}
	}
}

// LoadValidPasswordToken loads a valid password token entity that matches the user and token
// provided in path parameters
// If the token is invalid, the user will be redirected to the forgot password route
// This requires that the user owning the token is loaded in to context.
func LoadValidPasswordToken(authClient *services.AuthClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Extract the user parameter
			if c.Get(context.UserKey) == nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			}
			usr := c.Get(context.UserKey).(*ent.User)

			// Extract the token ID.
			tokenID, err := strconv.Atoi(c.Param("password_token"))
			if err != nil {
				return echo.NewHTTPError(http.StatusNotFound)
			}

			// Attempt to load a valid password token.
			token, err := authClient.GetValidPasswordToken(
				c,
				usr.ID,
				tokenID,
				c.Param("token"),
			)

			switch err.(type) {
			case nil:
				c.Set(context.PasswordTokenKey, token)
				return next(c)
			case services.InvalidPasswordTokenError:
				msg.Warning(c, "The link is either invalid or has expired. Please request a new one.")
				return c.Redirect(http.StatusFound, c.Echo().Reverse(routenames.ForgotPassword))
			default:
				return echo.NewHTTPError(
					http.StatusInternalServerError,
					fmt.Sprintf("error loading password token: %v", err),
				)
			}
		}
	}
}

// RequireAuthentication requires that the user be authenticated in order to proceed.
func RequireAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u == nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		return next(c)
	}
}

// RequireNoAuthentication requires that the user not be authenticated in order to proceed.
func RequireNoAuthentication(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u != nil {
			return c.Redirect(http.StatusSeeOther, c.Echo().Reverse(routenames.Dashboard))
		}

		return next(c)
	}
}

// RequireAdmin requires that the authenticated user be an admin in order to proceed.
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if u := c.Get(context.AuthenticatedUserKey); u != nil {
			if user, ok := u.(*ent.User); ok {
				if user.Admin {
					return next(c)
				}
			}
		}

		return echo.NewHTTPError(http.StatusUnauthorized)
	}
}

// RequirePaidUser requires that the authenticated user has either an active subscription
// or a successful payment intent in order to proceed.
func RequirePaidUser(db *ent.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First ensure user is authenticated
			u := c.Get(context.AuthenticatedUserKey)
			if u == nil {
				return echo.NewHTTPError(http.StatusUnauthorized)
			}

			user, ok := u.(*ent.User)
			if !ok {
				return echo.NewHTTPError(http.StatusInternalServerError, "invalid user type")
			}

			// Check if user has an active subscription
			hasActiveSubscription, err := db.Subscription.
				Query().
				Where(subscription.HasCustomerWith(
					paymentcustomer.HasUserWith(entuser.IDEQ(user.ID)),
				)).
				Where(subscription.StatusEQ(subscription.StatusActive)).
				Exist(c.Request().Context())
			
			if err != nil {
				log.Ctx(c).Warn(fmt.Sprintf("error checking subscription status: %v", err))
			}

			if hasActiveSubscription {
				return next(c)
			}

			// Check if user has a successful payment intent (one-time purchase)
			hasSuccessfulPayment, err := db.PaymentIntent.
				Query().
				Where(paymentintent.HasCustomerWith(
					paymentcustomer.HasUserWith(entuser.IDEQ(user.ID)),
				)).
				Where(paymentintent.StatusEQ(paymentintent.StatusSucceeded)).
				Exist(c.Request().Context())
			
			if err != nil {
				log.Ctx(c).Warn(fmt.Sprintf("error checking payment intent status: %v", err))
			}

			if hasSuccessfulPayment {
				return next(c)
			}

			// User doesn't have valid payment, redirect to products page
			msg.Warning(c, "Premium access required. Please purchase a product or subscribe to a plan.")
			return c.Redirect(http.StatusSeeOther, c.Echo().Reverse(routenames.Products))
		}
	}
}
