package httpmock

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ExpectedMultipartForm struct {
	Files  map[string]map[string][]byte
	Values map[string][]string
}

type call struct {
	method                string
	route                 string
	returnStatus          int
	returnBody            string
	returnContentLength   int64
	returnHeaders         map[string]string
	expectedBody          string
	expectedBodyFunc      func(string) bool
	expectedJSON          string
	expectedMultipartForm *ExpectedMultipartForm
	err                   error
	headers               map[string][]string
	queryParams           map[string][]string
}

type CallOption func(*call)

func ReturnError(err error) CallOption {
	return func(c *call) {
		c.err = err
	}
}

func ReturnStatus(status int) CallOption {
	return func(c *call) {
		c.returnStatus = status
	}
}

func ReturnBody(body string) CallOption {
	return func(c *call) {
		c.returnBody = body
	}
}

func ReturnHeader(name, value string) CallOption {
	return func(c *call) {
		if c.returnHeaders == nil {
			c.returnHeaders = make(map[string]string)
		}
		c.returnHeaders[name] = value
	}
}

func ReturnContentLength(contentLength int64) CallOption {
	return func(c *call) {
		c.returnContentLength = contentLength
	}
}

func ExpectBody(expectedBody string) CallOption {
	return func(c *call) {
		c.expectedBody = expectedBody
	}
}

func ExpectBodyFunc(f func(expectedBody string) bool) CallOption {
	return func(c *call) {
		c.expectedBodyFunc = f
	}
}

func ExpectJSON(expectedJSON string) CallOption {
	return func(c *call) {
		c.expectedJSON = expectedJSON
	}
}

func ExpectMultipartForm(expectedMultipartForm ExpectedMultipartForm) CallOption {
	return func(c *call) {
		c.expectedMultipartForm = &expectedMultipartForm
	}
}

func ExpectHeader(name string, values []string) CallOption {
	return func(c *call) {
		if c.headers == nil {
			c.headers = make(map[string][]string)
		}
		c.headers[name] = values
	}
}

func ExpectQueryParam(name, value string) CallOption {
	return func(c *call) {
		if c.queryParams == nil {
			c.queryParams = make(map[string][]string)
		}
		c.queryParams[name] = append(c.queryParams[name], value)
	}
}

type mockTransport struct {
	t     *testing.T
	index int
	calls []call
}

func (c *mockTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.t.Helper()
	if !assert.Less(c.t, c.index, len(c.calls)) {
		return nil, fmt.Errorf("unexpected route call #%d", c.index)
	}
	call := c.calls[c.index]
	if !assert.Equal(c.t, call.route, r.URL.Path) || !assert.Equal(c.t, call.method, r.Method) {
		return nil, fmt.Errorf("unexpected route call #%d on route %s %s, expected %s %s",
			c.index, r.Method, r.URL.Path, call.method, call.route)
	}

	if call.expectedJSON != "" {
		if !assert.NotNil(c.t, r.Body) {
			return nil, fmt.Errorf("expected body but received nothing")
		}
		data, err := io.ReadAll(r.Body)
		if !assert.NoError(c.t, err) {
			return nil, fmt.Errorf("io.ReadAll error: %w", err)
		}
		if !assert.JSONEq(c.t, call.expectedJSON, string(data)) {
			return nil, fmt.Errorf("expected body does not match received body")
		}
	}

	if call.expectedBody != "" || call.expectedBodyFunc != nil {
		if !assert.NotNil(c.t, r.Body) {
			return nil, fmt.Errorf("expected body but received nothing")
		}
		data, err := io.ReadAll(r.Body)
		if !assert.NoError(c.t, err) {
			return nil, fmt.Errorf("io.ReadAll error: %w", err)
		}

		if call.expectedBody != "" {
			if !assert.Equal(c.t, call.expectedBody, string(data)) {
				return nil, fmt.Errorf("expected body does not match received body")
			}
		}

		if call.expectedBodyFunc != nil {
			if !assert.True(c.t, call.expectedBodyFunc(string(data))) {
				return nil, fmt.Errorf("received body does not match expected body: %s", string(data))
			}
		}
	}

	if call.expectedMultipartForm != nil {
		err := r.ParseMultipartForm(r.ContentLength)
		if err != nil {
			return nil, fmt.Errorf("could not parse multipart form: %w", err)
		}
		for fieldName, field := range call.expectedMultipartForm.Files {
			got := make(map[string][]byte)
			for _, part := range r.MultipartForm.File[fieldName] {
				file, err := part.Open()
				if err != nil {
					return nil, fmt.Errorf("could not open form part %s in field %s: %w", part.Filename, fieldName, err)
				}

				content, err := io.ReadAll(file)
				if err != nil {
					return nil, fmt.Errorf("could not read form part %s in field %s: %w", part.Filename, fieldName, err)
				}

				got[part.Filename] = content
			}

			for filename, content := range field {
				gotContent, ok := got[filename]
				if !ok {
					return nil, fmt.Errorf("could not find form part %s in field %s", filename, fieldName)
				}
				if !bytes.Equal(content, gotContent) {
					return nil, fmt.Errorf("received content (%s) does not match expected content (%s) for part %s in field %s", gotContent, content, filename, fieldName)
				}
			}
		}

		for fieldName, field := range call.expectedMultipartForm.Values {
			for _, content := range field {
				var found bool
				for _, got := range r.MultipartForm.Value[fieldName] {
					if strings.TrimSpace(content) == strings.TrimSpace(got) {
						found = true
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("content (%s) not found in (%s) with field name %s", content, r.MultipartForm.Value[fieldName], fieldName)
				}
			}
		}
	}

	for name, values := range call.headers {
		requestValues, ok := r.Header[name]
		if !assert.True(c.t, ok) {
			return nil, fmt.Errorf("header %q not set in request", name)
		}
		if !assert.Equal(c.t, values, requestValues) {
			return nil, fmt.Errorf("header %q has bad values", name)
		}
	}

	requestQuery := r.URL.Query()
	for name, values := range call.queryParams {
		for _, value := range values {
			if !slices.Contains(requestQuery[name], value) {
				return nil, fmt.Errorf("query parameter %q: %q is missing (found: %q)", name, value, requestQuery[name])
			}
		}
	}

	c.index++
	if call.err != nil {
		return nil, call.err
	}
	resp := &http.Response{
		Status:        http.StatusText(call.returnStatus),
		StatusCode:    call.returnStatus,
		Body:          io.NopCloser(strings.NewReader(call.returnBody)),
		ContentLength: call.returnContentLength,
		Header:        make(http.Header),
	}
	for k, v := range call.returnHeaders {
		resp.Header.Set(k, v)
	}
	return resp, nil
}

type Client struct {
	http.Client
	tracker *mockTransport
}

func (c *Client) WithCall(method, route string, options ...CallOption) *Client {
	call := call{
		method: method,
		route:  route,
	}
	for _, option := range options {
		option(&call)
	}
	c.tracker.calls = append(c.tracker.calls, call)
	return c
}

func New(t *testing.T) *Client {
	transport := &mockTransport{
		t:     t,
		calls: make([]call, 0),
	}
	return &Client{
		tracker:   transport,
		Transport: transport,
	}
}
