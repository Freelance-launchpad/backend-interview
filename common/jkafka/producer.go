package jkafka

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	kafkatrace "github.com/DataDog/dd-trace-go/contrib/segmentio/kafka-go/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jhttp/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/internal"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/internal/registry"
	"github.com/kaptinlin/jsonschema"
	kafka "github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

type ProducerConfig struct {
	URL      string
	User     string
	Password string
	Topic    string `yaml:"topic"`

	SchemaRegistryConfig registry.Config `yaml:"schema_registry"`
}

func (pc *ProducerConfig) LoadFromEnvVars() {
	pc.URL = os.Getenv("KAFKA_URL")
	pc.User = os.Getenv("KAFKA_USER")
	pc.Password = os.Getenv("KAFKA_PASSWORD")
}

type Producer interface {
	RegisterSchema(ctx context.Context, schema string) error
	Write(ctx context.Context, message Message) error
}

type producer struct {
	kw internal.KafkaWriter

	registryClient internal.RegistryClient

	topic string

	schemaID   *int32
	jsonSchema *jsonschema.Schema
}

func NewProducer(config ProducerConfig) (Producer, error) {
	kafkaWriter := kafkatrace.WrapWriter(&kafka.Writer{
		Addr:        kafka.TCP(config.URL),
		Topic:       config.Topic,
		Balancer:    &kafka.Hash{},
		ErrorLogger: kafka.LoggerFunc(kafkaErrorLogger),
	})
	if config.User != "" {
		scram, err := scram.Mechanism(scram.SHA512, config.User, config.Password)
		if err != nil {
			return nil, fmt.Errorf("failed to create SCRAM mechanism: %w", err)
		}
		kafkaWriter.Transport = &kafka.Transport{
			SASL: scram,
		}
	}

	p := &producer{
		kw:    kafkaWriter,
		topic: config.Topic,
	}

	if config.SchemaRegistryConfig.URL.URL != nil {
		p.registryClient = registry.New(
			jhttp.NewClient(5*time.Second).
				WithTracing("schema-registry"),
			config.SchemaRegistryConfig,
		)
	}

	return p, nil
}

// RegisterSchema registers the producer's schema in the Schema Registry and stores the returned schema ID.
func (p *producer) RegisterSchema(ctx context.Context, schema string) (err error) {
	defer jerror.Wrap(&err)

	if p.registryClient == nil {
		return ErrRegistryClientNotConfigured
	}

	jsonSchema, err := jsonschema.NewCompiler().Compile([]byte(schema))
	if err != nil {
		return fmt.Errorf("failed to compile JSON schema: %w", err)
	}

	schemaID, err := p.registryClient.RegisterSchema(ctx, p.topic+"-value", schema)
	if err != nil {
		return err
	}

	p.schemaID = &schemaID
	p.jsonSchema = jsonSchema

	return nil
}

func (p *producer) Write(ctx context.Context, m Message) (err error) {
	key := m.Key()

	defer jerror.Wrap(&err, "with key", string(key))

	message, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if p.schemaID == nil {

		return p.kw.WriteMessages(ctx, kafka.Message{
			Key:   []byte(key),
			Value: message,
		})
	}

	if p.jsonSchema == nil {
		return ErrMissingJSONSchema
	}

	res := p.jsonSchema.ValidateJSON(message)
	if !res.IsValid() {
		detailedErrors, err := json.Marshal(res.DetailedErrors())
		if err != nil {
			detailedErrors = []byte(fmt.Sprintf("failed to marshal detailed errors: %v", err))
		}
		return fmt.Errorf("message does not conform to the registered schema: %+v: %w", string(detailedErrors), ErrInvalidMessage)
	}

	encoded, err := registry.Encode(*p.schemaID, message)
	if err != nil {
		return fmt.Errorf("failed to encode message: %w", err)
	}

	return p.kw.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: encoded,
	})
}
