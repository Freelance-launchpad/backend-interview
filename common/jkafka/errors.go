package jkafka

import "errors"

var (
	ErrInvalidMessage              = errors.New("invalid message")
	ErrRegistryClientNotConfigured = errors.New("schema registry client is not configured")
	ErrMissingJSONSchema           = errors.New("missing json schema")
)
