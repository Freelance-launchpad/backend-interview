package jhttp

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

// APIError is an error that can be used by clients to represent an error from the API.
// Duplicated from jhttp in order to keep the import clean.
type APIError struct {
	Err        error
	StatusCode int
}

func (e APIError) Error() string {
	return fmt.Sprintf("[%d]: %s", e.StatusCode, e.Err)
}

func (e APIError) Unwrap() error {
	return e.Err
}

func IsNotFoundAPIError(err error) bool {
	if apiErr, ok := errors.AsType[APIError](err); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

func IsHTTPAPIError(err error, status int) bool {
	if apiErr, ok := errors.AsType[APIError](err); ok {
		return apiErr.StatusCode == status
	}
	return false
}

// Doer is an interface used to easily abstract http clients
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	*http.Client
}

func NewClient(timeout time.Duration) *Client {
	client := &Client{
		Client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   timeout,
		},
	}
	return client.WithRequestID()
}

func (c *Client) WithRequestID() *Client {
	c.Transport = roundTripperRequestID{
		next: c.Transport,
	}

	return c
}

func (c *Client) WithBasicAuth(username, password string) *Client {
	c.Transport = roundTripperBasicAuth{
		next:     c.Transport,
		username: username,
		password: password,
	}
	return c
}

func (c *Client) WithBearerAuthentication(auth BearerAuthentication) *Client {
	c.Transport = roundTripperBearerAuthentication{
		next: c.Transport,
		auth: auth,
	}

	return c
}

func (c *Client) WithUserAgent(appName string, appVersion string) *Client {
	c.Transport = roundTripperUserAgent{
		next:      c.Transport,
		userAgent: fmt.Sprintf("%s/%s", appName, appVersion),
	}

	return c
}

func (c *Client) WithTracing(name string) *Client {
	c.Transport = roundTripperTracing{
		next:          c.Transport,
		operationName: fmt.Sprintf("%s.http.client", name),
	}

	return c
}

func (c *Client) WithRateLimiting(defaultLimit Limiter, customLimit ...PerRouteLimiter) *Client {
	defaultLimiter := rate.NewLimiter(defaultLimit.Limit, defaultLimit.Burst)

	customLimiter := map[string]*rate.Limiter{}
	for _, custom := range customLimit {
		customLimiter[custom.SubString] = rate.NewLimiter(custom.Limit, custom.Burst)
	}

	c.Transport = roundTripperRateLimiting{
		defaultLimiter: defaultLimiter,
		customLimiter:  customLimiter,
		next:           c.Transport,
	}

	return c
}

func (c *Client) WithRoundTripperFunc(fn RoundTripperFunc) *Client {
	c.Transport = roundTripperFn{
		fn:   fn,
		next: c.Transport,
	}
	return c
}
