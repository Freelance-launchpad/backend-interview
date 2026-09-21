package jhttp

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"golang.org/x/time/rate"
)

type CTXKey string

const (
	CTXKeyPriority     CTXKey = "priority"
	CTXKeyWaitDuration CTXKey = "waitDuration"
)

var _ http.RoundTripper = (*roundTripperRateLimiting)(nil)

type roundTripperRateLimiting struct {
	defaultLimiter *rate.Limiter
	customLimiter  map[string]*rate.Limiter
	next           http.RoundTripper
}

type Limiter struct {
	Limit rate.Limit
	Burst int
}

type PerRouteLimiter struct {
	SubString string
	Limiter
}

func (r roundTripperRateLimiting) RoundTrip(request *http.Request) (*http.Response, error) {
	ctx := request.Context()

	rateLimiter := r.defaultLimiter
	for substring, limiter := range r.customLimiter {
		if strings.Contains(request.URL.String(), substring) {
			rateLimiter = limiter
			break
		}
	}

	start := time.Now()
	if err := rateLimiter.Wait(ctx); err != nil {
		return nil, fmt.Errorf("rateLimiter.Wait failed: %w", err)
	}

	ctx = context.WithValue(ctx, CTXKeyWaitDuration, time.Since(start))
	request = request.WithContext(ctx)

	return r.next.RoundTrip(request)
}
