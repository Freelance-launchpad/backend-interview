package jhttp

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"strconv"

	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
)

var _ http.RoundTripper = (*roundTripperTracing)(nil)

type roundTripperTracing struct {
	next          http.RoundTripper
	operationName string
}

func (r roundTripperTracing) RoundTrip(request *http.Request) (*http.Response, error) {
	var (
		dnsSpan     *tracer.Span
		connectSpan *tracer.Span
		tlsSpan     *tracer.Span
	)

	span, ctx := tracer.StartSpanFromContext(request.Context(), r.operationName,
		tracer.SpanType(ext.SpanTypeHTTP),
	)
	defer span.Finish()

	trace := &httptrace.ClientTrace{
		DNSStart: func(dns httptrace.DNSStartInfo) {
			dnsSpan, _ = tracer.StartSpanFromContext(ctx, "dns.resolution",
				tracer.SpanType("dns"),
				tracer.ResourceName(dns.Host),
				tracer.ServiceName("dns"),
			)
		},
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			if dnsInfo.Err != nil {
				span.SetTag(ext.Error, dnsInfo.Err)
			}
			dnsSpan.SetTag("dns.addresses", dnsInfo.Addrs)
			dnsSpan.Finish()
		},

		ConnectStart: func(network, addr string) {
			connectSpan, _ = tracer.StartSpanFromContext(ctx, "http.connect",
				tracer.SpanType("connect"),
				tracer.ResourceName(network+"_"+addr),
			)
		},
		ConnectDone: func(_, _ string, err error) {
			connectSpan.Finish(tracer.WithError(err))
		},

		TLSHandshakeStart: func() {
			tlsSpan, _ = tracer.StartSpanFromContext(ctx, "http.tls",
				tracer.SpanType("tls_handshake"),
			)
		},
		TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
			tlsSpan.Finish(tracer.WithError(err))
		},
	}

	request = request.WithContext(httptrace.WithClientTrace(ctx, trace))
	err := tracer.Inject(span.Context(), tracer.HTTPHeadersCarrier(request.Header))
	if err != nil {
		span.SetTag("tracer_inject_error", err)
	}

	res, err := r.next.RoundTrip(request)
	if err != nil {
		span.SetTag("http.errors", err.Error())
		span.SetTag(ext.Error, err)
	} else {
		span.SetTag(ext.HTTPCode, strconv.Itoa(res.StatusCode))
		if res.StatusCode >= 500 {
			span.SetTag("http.errors", res.Status)
			span.SetTag(ext.Error, fmt.Errorf("%d: %s", res.StatusCode, http.StatusText(res.StatusCode)))
		}
	}

	return res, err
}
