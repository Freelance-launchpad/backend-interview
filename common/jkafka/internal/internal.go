package internal

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type RegistryClient interface {
	// RegisterSchema registers a JSON Schema under the given subject and returns the schema ID.
	// If the schema is already registered, the existing ID is returned.
	RegisterSchema(ctx context.Context, subject, schema string) (int32, error)

	// GetSchema fetches the raw schema string for a given schema ID.
	GetSchema(ctx context.Context, id int32) (string, error)
}

type KafkaWriter interface {
	WriteMessages(context.Context, ...kafka.Message) error
}
