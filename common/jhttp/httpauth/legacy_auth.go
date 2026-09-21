package httpauth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
)

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type authToken struct {
	token          string
	expirationDate time.Time
}

// LegacyConfiguration defines the authentication http client configuration.
type LegacyConfiguration struct {
	AuthURL      string `yaml:"auth_url"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Audience     string `yaml:"audience"`
	Scope        string `yaml:"scope"`
	GrantType    string `yaml:"grant_type"`
	AsForm       bool   `yaml:"as_form"`
}

var _ jhttp.BearerAuthentication = (*LegacyAuthentication)(nil)

// Deprecated: use Auth instead.
type LegacyAuthentication struct {
	mu    sync.Mutex
	conf  LegacyConfiguration
	token authToken
	doer  doer
}

type getAuthTokenRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience,omitempty"`
	Scope        string `json:"scope,omitempty"`
	GrantType    string `json:"grant_type"`
}

type getAuthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresInSec int    `json:"expires_in"`
}

func buildFormRequest(ctx context.Context, conf LegacyConfiguration) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	addStringField := func(name string, value string) error {
		fw, err := writer.CreateFormField(name)
		if err != nil {
			return fmt.Errorf("writer.CreateFormField error: %w", err)
		}
		_, err = io.Copy(fw, strings.NewReader(value))
		if err != nil {
			return fmt.Errorf("io.Copy error: %w", err)
		}

		return nil
	}

	if err := addStringField("client_id", conf.ClientID); err != nil {
		return nil, fmt.Errorf("addStringField %q error: %w", "client_id", err)
	}
	if err := addStringField("client_secret", conf.ClientSecret); err != nil {
		return nil, fmt.Errorf("addStringField %q error: %w", "client_secret", err)
	}
	if err := addStringField("audience", conf.Audience); err != nil {
		return nil, fmt.Errorf("addStringField %q error: %w", "audience", err)
	}
	if err := addStringField("scope", conf.Scope); err != nil {
		return nil, fmt.Errorf("addStringField %q error: %w", "scope", err)
	}
	if err := addStringField("grant_type", conf.GrantType); err != nil {
		return nil, fmt.Errorf("addStringField %q error: %w", "grant_type", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("writer.Close error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.AuthURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext error: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func buildJSONRequest(ctx context.Context, conf LegacyConfiguration) (*http.Request, error) {
	data := getAuthTokenRequest{
		ClientID:     conf.ClientID,
		ClientSecret: conf.ClientSecret,
		Audience:     conf.Audience,
		GrantType:    conf.GrantType,
		Scope:        conf.Scope,
	}

	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, conf.AuthURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("http.NewRequestWithContext error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func refreshToken(ctx context.Context, httpClient doer, conf LegacyConfiguration) (authToken, error) {
	errMsg := "Authentication.refreshToken has failed"

	var req *http.Request
	var err error
	if conf.AsForm {
		req, err = buildFormRequest(ctx, conf)
	} else {
		req, err = buildJSONRequest(ctx, conf)
	}
	if err != nil {
		return authToken{}, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return authToken{}, fmt.Errorf("%s: http.Client.Do error: %w", errMsg, err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			responseBody = []byte("error reading response body")
		}
		return authToken{}, fmt.Errorf("%s: bad status, expected: %d, received: %d - with body %s", errMsg, http.StatusOK, resp.StatusCode, string(responseBody))
	}

	body := getAuthTokenResponse{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return authToken{}, fmt.Errorf("json.NewDecoder().Decode error: %w", err)
	}

	return authToken{
		token: body.AccessToken,
		expirationDate: time.Now().
			Add(time.Duration(body.ExpiresInSec) * time.Second).
			Add(-30 * time.Second),
	}, nil
}

func (a *LegacyAuthentication) GetBearer(ctx context.Context) (string, error) {
	if a.token.expirationDate.After(time.Now()) {
		return a.token.token, nil
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	if a.token.expirationDate.Before(time.Now()) {
		token, err := refreshToken(ctx, a.doer, a.conf)
		if err != nil {
			return "", fmt.Errorf("getBearer has failed: refreshToken err: %w", err)
		}
		a.token = token
	}

	return a.token.token, nil
}

// Deprecated: use NewAuth instead.
func New(conf LegacyConfiguration, doer doer) *LegacyAuthentication {
	return &LegacyAuthentication{
		conf: conf,
		doer: doer,
	}
}
