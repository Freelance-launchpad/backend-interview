package jhttp

import "net/http"

var _ http.RoundTripper = (*roundTripperFn)(nil)

type RoundTripperFunc func(r *http.Request, next http.RoundTripper) (*http.Response, error)

type roundTripperFn struct {
	fn   RoundTripperFunc
	next http.RoundTripper
}

func (rt roundTripperFn) RoundTrip(req *http.Request) (*http.Response, error) {
	return rt.fn(req, rt.next)
}
