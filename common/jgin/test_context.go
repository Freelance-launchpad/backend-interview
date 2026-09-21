package jgin

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/Freelance-launchpad/backend-interview/common/jentity"
	pkgusers "github.com/Freelance-launchpad/backend-interview/services/users/pkg"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ContextOption configures the gin context in order to create a good test environment
type ContextOption func(*gin.Context)

// CreateTestContext a gin context in test mode then fill the context thanks to the optional functions
func CreateTestContext(method string, options ...ContextOption) (*httptest.ResponseRecorder, *gin.Context) {
	// Set gin in test mode
	gin.SetMode(gin.TestMode)

	// Create context
	w := httptest.NewRecorder()
	ctx, engine := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
		Method: method,
	}
	engine.TrustedPlatform = "X-header-test-remote-addr"

	for _, opt := range options {
		opt(ctx)
	}

	return w, ctx
}

func WithUserID(userID string) ContextOption {
	return func(ctx *gin.Context) {
		actualCTX := ctx.Request.Context()
		newCTX := ContextWithUserID(actualCTX, userID)
		ctx.Request = ctx.Request.WithContext(newCTX)
	}
}

func WithEmailVerified(emailVerified bool) ContextOption {
	return func(ctx *gin.Context) {
		actualCTX := ctx.Request.Context()
		newCTX := ContextWithEmailVerified(actualCTX, emailVerified)
		ctx.Request = ctx.Request.WithContext(newCTX)
	}
}

// Deprecated
func WithUserEmail(userEmail string) ContextOption {
	return func(ctx *gin.Context) {
		actualCTX := ctx.Request.Context()
		newCTX := ContextWithUserEmail(actualCTX, userEmail)
		ctx.Request = ctx.Request.WithContext(newCTX)
	}
}

func WithJSONBody(input any) ContextOption {
	return func(ctx *gin.Context) {
		bodyBytes, err := json.Marshal(input)
		if err != nil {
			panic(err)
		}

		ctx.Request.ContentLength = int64(len(bodyBytes))
		ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
}

func WithBody(input io.ReadCloser) ContextOption {
	return func(ctx *gin.Context) {
		var buf bytes.Buffer
		if _, err := buf.ReadFrom(input); err != nil {
			panic(err)
		}

		ctx.Request.ContentLength = int64(buf.Len())
		ctx.Request.Body = io.NopCloser(&buf)
	}
}

func WithParams(key, value string) ContextOption {
	return func(ctx *gin.Context) {
		ctx.AddParam(key, value)
	}
}

func WithQueryParams(key string, values ...string) ContextOption {
	return func(ctx *gin.Context) {
		queryParams := ctx.Request.URL.Query()
		for _, value := range values {
			queryParams.Add(key, value)
		}
		ctx.Request.URL.RawQuery = queryParams.Encode()
	}
}

func WithOfferID(offerID uuid.UUID) ContextOption {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(ContextWithOfferID(ctx.Request.Context(), offerID))
	}
}

func WithOfferType(offerType jentity.OfferType) ContextOption {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(ContextWithOfferType(ctx.Request.Context(), offerType))
	}
}

func WithEntity(entity jentity.Entity) ContextOption {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(ContextWithEntity(ctx.Request.Context(), entity))
	}
}

func WithOfferContract(oc pkgusers.OfferWithContract) ContextOption {
	return func(ctx *gin.Context) {
		ctx.Request = ctx.Request.WithContext(ContextWithOfferContract(ctx.Request.Context(), oc))
	}
}

func WithHeader(header string, values ...string) ContextOption {
	return func(ctx *gin.Context) {
		for _, value := range values {
			ctx.Request.Header.Add(header, value)
		}
	}
}

func WithRemoteInfo(info RemoteInformation) ContextOption {
	return func(c *gin.Context) {
		if info.RemoteUserID != nil {
			c.Request = c.Request.WithContext(ContextWithUserID(c.Request.Context(), *info.RemoteUserID))
		}
		if info.RemoteAddress != nil {
			c.Request.Header.Add("X-header-test-remote-addr", *info.RemoteAddress)
		}
		if info.RemoteUserAgent != nil {
			c.Request.Header.Add("User-Agent", *info.RemoteUserAgent)
		}
	}
}
