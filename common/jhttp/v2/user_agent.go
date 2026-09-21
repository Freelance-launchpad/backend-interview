package jhttp

import (
	"net/http"
)

var _ http.RoundTripper = (*roundTripperUserAgent)(nil)

type roundTripperUserAgent struct {
	next      http.RoundTripper
	userAgent string
}

func (r roundTripperUserAgent) RoundTrip(request *http.Request) (*http.Response, error) {
	request.Header.Add("User-Agent", r.userAgent)

	return r.next.RoundTrip(request)
}
