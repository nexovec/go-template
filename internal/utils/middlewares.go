package utils

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/gofiber/fiber/v3"
)

// This verifies that the incoming http Connect-Type header matches contentType
func BuildMwEnsureMIME(allowedTypes ...string) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Compare the Content-Type header with the expected value
		contentType := c.Get(fiber.HeaderContentType)
		print(allowedTypes)
		for _, allowedType := range allowedTypes {
			if strings.HasPrefix(contentType, allowedType) {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusBadRequest).SendString("Invalid Content-Type")
	}
}

func MwDefaultResponseMIME(contentType string) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, contentType)
		return c.Next()
	}
}

func GuestOnlyMwBuilder(redirectOnErrorURL string) fiber.Handler {
	return func(c fiber.Ctx) error {
		err := ParseClaimsIdemp(c)
		if err == nil && VerifyUser(c) == nil {
			return c.Redirect().To(redirectOnErrorURL)
		} else if err != nil && err != ErrNoCookie {
			slog.Error("login", "Error parsing claims:", err)
			WipeCookie(c, "Bearer-Token")
			return c.Redirect().To(redirectOnErrorURL)
		}
		return c.Next()
	}
}

type HandlerWithErrorPropagation func(fiber.Ctx, error) error

func LoggedInOnlyJWTMiddlewareBuilder(onFail HandlerWithErrorPropagation) fiber.Handler {
	if onFail == nil {
		onFail = func(c fiber.Ctx, err error) error {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
	}
	return func(c fiber.Ctx) error {
		parseErr := ParseClaimsIdemp(c)
		verifyErr := VerifyUser(c)
		err := errors.Join(parseErr, verifyErr)
		if err != nil {
			return onFail(c, err)
		}
		return c.Next()
	}
}

// This middleware is used to redirect to the top-level parent view if request isn't an AJAX.
// This is meant for use with view subroutes that don't make sense as a standalone page.
// This is specifically only usable for routes under the /v2/app/view router right now. // TODO: this is trivial to extend with a parameter
func NonAjaxRedirectToParentViewMw(c fiber.Ctx) error {
	// TODO: implement this
	return c.Next()
}
