package jgin

import (
	"crypto/ed25519"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Auth struct {
	jwtPublicKey ed25519.PublicKey
	bannedJTIs   []string
	legacyAuth   *LegacyAuth
	skipAuth     bool
}

func NewAuth(jwtPublicKey ed25519.PublicKey) *Auth {
	return &Auth{
		jwtPublicKey: jwtPublicKey,
	}
}

func (a *Auth) WithBannedJTIs(bannedJTIs []string) *Auth {
	a.bannedJTIs = bannedJTIs
	return a
}

func (a *Auth) WithFallback(legacyAuth *LegacyAuth) *Auth {
	a.legacyAuth = legacyAuth
	return a
}

func (a *Auth) WithSkipAuth(skipAuth bool) *Auth {
	a.skipAuth = skipAuth
	return a
}

func (a *Auth) authenticate(c *gin.Context) error { return nil }

func (a *Auth) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		if a.skipAuth {
			c.Request = c.Request.WithContext(ContextWithPermissions(c.Request.Context(), []string{"*"}))
			return
		}

		if err := a.authenticate(c); err != nil {
			if !errors.Is(err, jwt.ErrTokenSignatureInvalid) || a.legacyAuth == nil {
				_ = c.AbortWithError(http.StatusUnauthorized, err)
				return
			}

			// Fallback to Auth0 auth middleware
			a.legacyAuth.Authenticate()(c)
		}
	}
}
