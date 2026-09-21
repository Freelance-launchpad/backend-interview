package types

import (
	"github.com/Freelance-launchpad/backend-interview/common/jkafka"
	"github.com/google/uuid"
)

type StoredEvent struct {
	EventID   uuid.UUID
	Topic     string
	ID        string
	EventType jkafka.Type
	State     string
}
