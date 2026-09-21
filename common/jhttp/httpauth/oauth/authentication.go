package oauth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
)

type Token struct {
	Name         string
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

type Configuration struct {
	AuthURL      string `yaml:"auth_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"-"`
	TokenName    string `yaml:"token_name"`
}

type OAuth struct {
	conf    Configuration
	client  jhttp.Doer
	store   Store
	clock   jclock.Clock
	token   Token
	muToken sync.RWMutex
}

func NewOAuth(store Store, conf Configuration) (jhttp.BearerAuthentication, error) {
	auth := &OAuth{
		conf:   conf,
		client: jhttp.NewClient(10 * time.Second).WithTracing("authentication"),
		store:  store,
		clock:  jclock.RealClock{},
	}

	if _, err := auth.getToken(context.Background()); err != nil {
		return nil, fmt.Errorf("unable to init, name: %s: %w", auth.conf.TokenName, err)
	}

	return auth, nil
}

type reqRefreshToken struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RefreshToken string   `json:"refresh_token"`
	GrantType    string   `json:"grant_type"`
	Scopes       []string `json:"scopes,omitempty"`
}

type respRefreshToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (s *OAuth) refreshToken(ctx context.Context, token Token) (Token, error) {
	data, err := json.Marshal(reqRefreshToken{
		ClientID:     s.conf.ClientID,
		ClientSecret: s.conf.ClientSecret,
		RefreshToken: token.RefreshToken,
		GrantType:    "refresh_token",
	})
	if err != nil {
		return Token{}, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.conf.AuthURL, bytes.NewReader(data))
	if err != nil {
		return Token{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return Token{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = []byte(err.Error())
		}
		return Token{}, jhttp.APIError{
			Err:        fmt.Errorf("bad status, expected: %d, received: %d - with body %s", http.StatusOK, resp.StatusCode, string(responseBody)),
			StatusCode: resp.StatusCode,
		}
	}

	var body respRefreshToken
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return Token{}, err
	}

	return Token{
		Name:         token.Name,
		AccessToken:  body.AccessToken,
		RefreshToken: body.RefreshToken,
		ExpiresAt: s.clock.Now().UTC().
			Add(time.Duration(body.ExpiresIn) * time.Second).
			Add(-1 * time.Minute),
	}, nil
}

func (s *OAuth) getToken(ctx context.Context) (_ Token, err error) {
	defer jerror.Wrap(&err)

	s.muToken.RLock()
	if s.token.ExpiresAt.After(s.clock.Now()) {
		defer s.muToken.RUnlock()

		slog.DebugContext(ctx, "Got token from in memory variable", slog.String("name", s.conf.TokenName))
		return s.token, nil
	}
	s.muToken.RUnlock()

	s.muToken.Lock()
	defer s.muToken.Unlock()

	if s.token.ExpiresAt.Before(s.clock.Now()) {
		slog.DebugContext(ctx, "Refresh token with database value", slog.String("name", s.conf.TokenName))

		if s.token, err = s.store.GetToken(ctx, nil, s.conf.TokenName); err != nil {
			return Token{}, err
		}
	}

	if s.token.ExpiresAt.After(s.clock.Now()) {
		return s.token, nil
	}
	slog.DebugContext(ctx, "Refresh token with database value (lock row)", slog.String("name", s.conf.TokenName))

	tx, err := s.store.BeginLocked(ctx)
	if err != nil {
		return Token{}, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	s.token, err = s.store.GetToken(ctx, tx, s.conf.TokenName)
	if err != nil && !errors.Is(err, jerror.NotFoundError) {
		return Token{}, err
	}

	if s.token.ExpiresAt.Before(s.clock.Now()) {
		slog.DebugContext(ctx, "Refresh token with authentication service", slog.String("name", s.conf.TokenName))

		s.token, err = s.refreshToken(ctx, s.token)
		if err != nil {
			return Token{}, err
		}

		if err = s.store.UpdateToken(ctx, tx, s.token); err != nil {
			return Token{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return Token{}, err
	}

	return s.token, nil
}

func (s *OAuth) GetBearer(ctx context.Context) (bearer string, err error) {
	defer jerror.Wrap(&err, "with token_name", s.conf.TokenName)

	token, err := s.getToken(ctx)
	if err != nil {
		return "", err
	}

	return token.AccessToken, err
}
