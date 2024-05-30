package utils

import (
	"errors"
	"fmt"
	"models"
	"os"

	// "log"

	"time"

	_ "crypto/sha256"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	TokenId  uint      `json:"token_id"`
	UserId   int32     `json:"user_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	IsAdmin  bool      `json:"is_admin"`
	Ext      time.Time `json:"ext"`
	jwt.RegisteredClaims
}

var ENV_JWT_SIGNING_KEY string

var ErrNoCookie = errors.New("no cookie")

func (c *TokenClaims) PopulateFromRequestToken(ctx fiber.Ctx) error {
	user, ok := ctx.Locals("user_claims").(TokenClaims)
	if !ok {
		return errors.New("invalid user in context")
	}
	*c = user
	return c.Validate()
}
func (c *TokenClaims) Validate() error {
	// validation
	if c.TokenId == 0 || c.UserId == 0 || c.Email == "" {
		// TODO: this should not happen in prod and needs to be reported
		return errors.New("some token claims missing")
	}
	return nil
}

func FlushJWTCookieMw(c fiber.Ctx) error {
	// TODO: Inspect all places where this is used and invalidate the cache of logged in users if applicable
	WipeCookie(c, "Bearer-Token")
	return c.Next()
}

// sets c.Locals("user_claims") to the JWT's claims if it's not set already.
func ParseClaimsIdemp(c fiber.Ctx) error {
	local := c.Locals("user_claims")
	if local != nil {
		return nil
	}
	cookie := c.Cookies("Bearer-Token")
	if cookie == "" {
		return ErrNoCookie
	}
	var tokenClaims TokenClaims
	decodedToken, err := jwt.ParseWithClaims(cookie, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		// TODO: investigate the security implications of using this or other methods
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(ENV_JWT_SIGNING_KEY), nil
	})

	if err != nil {
		WipeCookie(c, "Bearer-Token")
		return err
	}
	// TODO: notify the user that they've been logged out because of an invalid token
	if !decodedToken.Valid {
		WipeCookie(c, "Bearer-Token")
		return errors.New("invalid token claims")
	}
	c.Locals("user_claims", tokenClaims)
	return nil
}

func VerifyUser(c fiber.Ctx) error {
	// NOTE: not necessary, if it fixes your issue, it means you're not calling function to ensure parsed claims properly
	err := ParseClaimsIdemp(c)
	if err != nil {
		return err
	}
	claims := TokenClaims{}
	err = claims.PopulateFromRequestToken(c)
	if err != nil {
		WipeCookie(c, "Bearer-Token")
		return errors.New("invalid JWT token claims")
	}
	if os.Getenv("DEPLOYMENT") == "debug" {
		user, err := models.New(models.Pool).GetUserById(c.Context(), claims.UserId)
		if err != nil {
			WipeCookie(c, "Bearer-Token")
			return errors.New("user does not exist")
		}
		// check if user is invalid
		if user.ID == 0 {
			WipeCookie(c, "Bearer-Token")
			return errors.New("user id is not valid")
		}
	}
	return nil
}

func GetTokenClaims(c fiber.Ctx) (TokenClaims, bool) {
	claims := c.Locals("user_claims")
	if claims == nil {
		return TokenClaims{}, false
	}
	claimsActual, ok := claims.(TokenClaims)
	return claimsActual, ok
}

func WipeCookie(c fiber.Ctx, cookieName string) {
	c.Cookie(&fiber.Cookie{Name: cookieName, Expires: time.Now().Add(-time.Second * 3600), Value: ""})
}
