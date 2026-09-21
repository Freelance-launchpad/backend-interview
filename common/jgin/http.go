package jgin

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const bearerPrefix = "Bearer "
const subSuffix = "@clients"

func (a *LegacyAuth) getTokenStringFromRequest(req *http.Request) string {
	authToken := req.Header.Get("Authorization")
	bearLen := len(bearerPrefix)
	if len(authToken) < bearLen || authToken[:bearLen] != bearerPrefix {
		return ""
	}
	return authToken[bearLen:]
}

// Authenticate returns an HTTP HandlerFunc authenticating a user.
func (a *LegacyAuth) Authenticate() gin.HandlerFunc {
	if a.skipAuth {
		return setTestUserData
	}

	return func(c *gin.Context) {
		tokenString := a.getTokenStringFromRequest(c.Request)
		if tokenString == "" {
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized"))
			return
		}

		token, err := a.ParseJWT(tokenString)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized: %w", err))
			return
		}

		a.setUserData(c, token)
	}
}

func (a *LegacyAuth) setUserData(c *gin.Context, token *jwt.Token) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized: invalid claims type"))
		return
	}

	sub, ok := claims["sub"].(string)
	if !ok {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized: missing subject"))
		return
	}

	ctx := c.Request.Context()

	if subParts := strings.Split(sub, "|"); len(subParts) == 2 {
		sub = subParts[1]
	} else if strings.HasSuffix(sub, subSuffix) {
		ctx = ServiceContext(ctx)
		sub = sub[:len(sub)-len(subSuffix)]
	} else {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized: invalid subject format"))
		return
	}

	ctx = ContextWithUserID(ctx, sub)

	email, ok := claims[fmt.Sprintf("%s/email", a.customFieldPrefix)].(string)
	if ok {
		ctx = ContextWithUserEmail(ctx, email)
	}

	ctx = ContextWithPermissions(ctx, auth0Scopes(claims))

	c.Request = c.Request.WithContext(ctx)
}

func setTestUserData(c *gin.Context) {
	userID := c.GetHeader("user_id_test")
	ctx := ContextWithUserID(c.Request.Context(), userID)
	ctx = ContextWithUserEmail(ctx, "foo@test.com")
	c.Request = c.Request.WithContext(ctx)
}
