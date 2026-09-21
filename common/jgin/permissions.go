package jgin

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// RequirePermissions *must* be called after Authenticate.
func RequirePermissions(requiredPermissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		permissions := GetPermissionsFromContext(ctx)

		if hasPermissions(permissions, requiredPermissions) {
			return
		}

		slog.DebugContext(ctx, "missing permissions",
			slog.Any("requiredPermissions", requiredPermissions),
			slog.Any("permissions", permissions),
		)

		c.AbortWithStatus(http.StatusForbidden)
	}
}

func HasPermissions(ctx context.Context, requiredPermissions ...string) bool {
	return hasPermissions(GetPermissionsFromContext(ctx), requiredPermissions)
}

// hasPermissions returns true if one of the requiredPermissions is not found in permissions.
func hasPermissions(permissions, requiredPermissions []string) bool {
	for _, requiredPermission := range requiredPermissions {
		if !wildcardStrategy(permissions, requiredPermission) {
			return false
		}
	}
	return true
}

// auth0Scopes returns a list of scopes found in claims.
func auth0Scopes(claims jwt.MapClaims) []string {
	scopes := []string{}

	if serviceScopesStr, success := claims["scope"].(string); success {
		serviceScopes := strings.SplitSeq(serviceScopesStr, " ")
		for scope := range serviceScopes {
			scopeParts := strings.Split(scope, "/")
			scopes = append(scopes, scopeParts[len(scopeParts)-1])
		}
	}

	if userScopes, success := claims["permissions"].([]any); success {
		for _, elem := range userScopes {
			if scope, success := elem.(string); success {
				scopes = append(scopes, scope)
			}
		}
	}

	return scopes
}

// HasPermission returns an HTTP HandlerFunc validating that a request has the correct right.
//
// Deprecated: use HasPermissions instead.
func (a *LegacyAuth) HasPermission(scope string) gin.HandlerFunc {
	return a.HasPermissions(scope)
}

// HasPermissions returns an HTTP HandlerFunc validating that a request has the correct rights.
func (a *LegacyAuth) HasPermissions(scopes ...string) gin.HandlerFunc {
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

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("not Authorized: invalid claims type"))
			return
		}

		if !hasPermissions(auth0Scopes(claims), scopes) {
			slog.InfoContext(c.Request.Context(), "Scope is missing from token", slog.Any("scope", scopes), slog.Any("subject", claims["sub"]))
			_ = c.AbortWithError(http.StatusForbidden, fmt.Errorf("forbidden"))
			return
		}

		a.setUserData(c, token)
	}
}

func wildcardStrategy(matchers []string, needle string) bool {
	needleParts := strings.Split(needle, ".")
	for _, matcher := range matchers {
		matcherParts := strings.Split(matcher, ".")

		if len(matcherParts) > len(needleParts) {
			continue
		}

		var noteq bool
		for k, c := range strings.Split(matcher, ".") {
			// this is the last item and the lengths are different
			if k == len(matcherParts)-1 && len(matcherParts) != len(needleParts) {
				if c != "*" {
					noteq = true
					break
				}
			}

			if c == "*" && len(needleParts[k]) > 0 {
				// pass because this satisfies the requirements
				continue
			} else if c != needleParts[k] {
				noteq = true
				break
			}
		}

		if !noteq {
			return true
		}
	}

	return false
}
