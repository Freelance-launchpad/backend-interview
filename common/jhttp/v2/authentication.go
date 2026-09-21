package jhttp

import (
	"context"
	"fmt"
	"net/http"
)

type BearerAuthentication interface {
	GetBearer(ctx context.Context) (string, error)
}

type StaticBearerToken struct {
	Token string
}

//nolint:unparam // error is always nil but required for the "BearerAuthentication" interface
func (s *StaticBearerToken) GetBearer(_ context.Context) (string, error) {
	return s.Token, nil
}

var _ http.RoundTripper = (*roundTripperBearerAuthentication)(nil)

type roundTripperBearerAuthentication struct {
	next http.RoundTripper
	auth BearerAuthentication
}

func (r roundTripperBearerAuthentication) RoundTrip(request *http.Request) (*http.Response, error) {
	accessToken, err := r.auth.GetBearer(request.Context())
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	return r.next.RoundTrip(request)
}

var _ http.RoundTripper = (*roundTripperBasicAuth)(nil)

type roundTripperBasicAuth struct {
	next     http.RoundTripper
	username string
	password string
}

func (r roundTripperBasicAuth) RoundTrip(request *http.Request) (*http.Response, error) {
	request.SetBasicAuth(r.username, r.password)
	return r.next.RoundTrip(request)
}
