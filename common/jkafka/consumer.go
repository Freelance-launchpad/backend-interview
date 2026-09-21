package jkafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	kafkatrace "github.com/DataDog/dd-trace-go/contrib/segmentio/kafka-go/v2"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/ext"
	"github.com/DataDog/dd-trace-go/v2/ddtrace/tracer"
	"github.com/Freelance-launchpad/backend-interview/common/japplication"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/internal"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/internal/registry"
	"github.com/cenkalti/backoff/v7"
	"github.com/kaptinlin/jsonschema"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type ConsumerConfig struct {
	URL      string
	User     string
	Password string
	Topic    string `yaml:"topic"`
	GroupID  string `yaml:"group_id"`

	MaxWait *time.Duration `yaml:"max_wait"`

	SchemaRegistryConfig registry.Config `yaml:"schema_registry"`
}

func (cc *ConsumerConfig) LoadFromEnvVars() {
	cc.URL = os.Getenv("KAFKA_URL")
	cc.User = os.Getenv("KAFKA_USER")
	cc.Password = os.Getenv("KAFKA_PASSWORD")
}

type consumer[T any] struct {
	reader  *kafkatrace.Reader
	handler Handler[T]
	done    chan struct{}

	registryClient internal.RegistryClient

	topic string

	schemasCache map[int32]*jsonschema.Schema
}

type Handler[T any] func(ctx context.Context, msg T) error

func NewConsumer[T any](config ConsumerConfig, handler Handler[T]) (japplication.Service, error) {
	kafkaConfig := kafka.ReaderConfig{
		Brokers:     []string{config.URL},
		GroupID:     config.GroupID,
		Topic:       config.Topic,
		MaxBytes:    1_000_000,
		ErrorLogger: kafka.LoggerFunc(kafkaErrorLogger),
	}
	if config.MaxWait != nil {
		kafkaConfig.MaxWait = *config.MaxWait
	}

	if config.User != "" {
		scram, err := scram.Mechanism(scram.SHA512, config.User, config.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to create SCRAM mechanism: %w", err)
		}
		kafkaConfig.Dialer = &kafka.Dialer{
			SASLMechanism: scram,
		}
	}

	c := &consumer[T]{
		reader:       kafkatrace.WrapReader(kafka.NewReader(kafkaConfig)),
		handler:      handler,
		done:         make(chan struct{}),
		topic:        config.Topic,
		schemasCache: map[int32]*jsonschema.Schema{},
	}

	if config.SchemaRegistryConfig.URL.URL != nil {
		c.registryClient = registry.New(
			jhttp.NewClient(5*time.Second).
				WithTracing("schema-registry"),
			config.SchemaRegistryConfig,
		)
	}

	return c, nil
}

func (c *consumer[_]) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case <-c.done:
			cancel()
		case <-ctx.Done():
		}
	}()

	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) ||
				errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		if _, err := backoff.Retry(ctx, func() (string, error) {
			if err := c.handle(ctx, msg); err != nil {
				return "", err
			}
			return "ok", nil
		},
			backoff.WithBackOff(backoff.NewExponentialBackOff()),
			backoff.WithMaxElapsedTime(0), // retry until context is done
			backoff.WithNotify(func(err error, d time.Duration) {
				slog.ErrorContext(ctx, "error handling kafka message",
					slog.String("topic", msg.Topic),
					slog.String("duration", d.String()),
					slog.Int("partition", msg.Partition),
					slog.Int64("offset", msg.Offset),
					slog.Any("error", err),
				)
			}),
		); err != nil {
			return err
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			slog.ErrorContext(ctx, "failed to commit kafka message",
				slog.String("topic", msg.Topic),
				slog.Int("partition", msg.Partition),
				slog.Int64("offset", msg.Offset),
				slog.Any("error", err),
			)
		}
	}
}

func (c *consumer[T]) handle(ctx context.Context, msg kafka.Message) (err error) {
	span, ctx := tracer.StartSpanFromContext(ctx, "kafka.Message",
		tracer.SpanType(ext.SpanTypeMessageConsumer),
		tracer.ResourceName(c.topic),
	)
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic error recovered: %v", r)
		}

		if err != nil {
			span.SetTag(ext.Error, err)
		}
		span.Finish()
	}()

	if c.registryClient != nil {
		schemaID, payload, err := registry.Decode(msg.Value)
		if err != nil {
			return fmt.Errorf("failed to decode message: %w", err)
		}

		schema, ok := c.schemasCache[schemaID]
		if !ok {
			schemaStr, err := c.registryClient.GetSchema(ctx, schemaID)
			if err != nil {
				return err
			}
			schema, err = jsonschema.NewCompiler().Compile([]byte(schemaStr))
			if err != nil {
				return fmt.Errorf("failed to compile JSON schema: %w", err)
			}
			c.schemasCache[schemaID] = schema
			slog.InfoContext(ctx, "fetched and compiled JSON schema from registry",
				slog.Int64("schemaID", int64(schemaID)),
			)
		}

		res := schema.ValidateJSON(payload)
		if !res.IsValid() {
			detailedErrors, err := json.Marshal(res.DetailedErrors())
			if err != nil {
				detailedErrors = []byte(fmt.Sprintf("failed to marshal detailed errors: %v", err))
			}
			return fmt.Errorf("message does not conform to schema: %s: %w", string(detailedErrors), ErrInvalidMessage)
		}

		msg.Value = payload
	}

	var data T
	if err := json.Unmarshal(msg.Value, &data); err != nil {
		return fmt.Errorf("failed to unmarshal message value: %w", err)
	}

	if err := c.handler(ctx, data); err != nil {
		return fmt.Errorf("handler error: %w", err)
	}

	return nil
}

func (c *consumer[_]) Stop(_ context.Context) error {
	close(c.done)
	return c.reader.Close()
}
