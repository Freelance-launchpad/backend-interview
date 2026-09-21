package jgin

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lestrrat-go/jwx/v3/jwk"
)

// LegacyAuth hold all information for validating a JWT token
type LegacyAuth struct {
	jwk               jwk.Set
	iss               string
	customFieldPrefix string
	skipAuth          bool
}

// NewLegacyAuth returns a new Auth object.
func NewLegacyAuth(config *Config) (*LegacyAuth, error) {
	if config.SkipAuth {
		return &LegacyAuth{
			skipAuth: true,
		}, nil
	}

	a := &LegacyAuth{
		iss:               config.ISS,
		customFieldPrefix: config.CustomFieldPrefix,
	}

	keySet, err := jwk.Fetch(context.Background(), config.JWKURL)
	if err != nil {
		return nil, err
	}

	a.jwk = keySet

	return a, nil
}

// ParseJWT parse a token string and returns a full jwt.Token object.
func (a *LegacyAuth) ParseJWT(tokenString string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, errors.New("kid header not found")
		}

		key, exist := a.jwk.LookupKeyID(kid)
		if !exist {
			return nil, fmt.Errorf("key %v not found", kid)
		}

		var raw any
		if err := jwk.Export(key, &raw); err != nil {
			return nil, err
		}

		return raw, nil
	})
	if err != nil {
		return nil, err
	}

	if !parsedToken.Valid {
		return nil, errors.New("invalid token")
	}

	iss, ok := parsedToken.Claims.(jwt.MapClaims)["iss"].(string)
	if !ok || iss != a.iss {
		return nil, errors.New("invalid iss field in token claims")
	}

	return parsedToken, nil
}
