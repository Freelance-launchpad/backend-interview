package jhttp

import (
	"net/http"

	"github.com/Freelance-launchpad/backend-interview/common/jgin/middlewares"
)

var _ http.RoundTripper = (*roundTripperRequestID)(nil)

type roundTripperRequestID struct {
	next http.RoundTripper
}

func (r roundTripperRequestID) RoundTrip(request *http.Request) (*http.Response, error) {
	requestID := request.Context().Value(middlewares.CtxRequestID)
	if id, ok := requestID.(string); ok {
		request.Header.Add("X-Request-Id", id)
	}

	return r.next.RoundTrip(request)
}
