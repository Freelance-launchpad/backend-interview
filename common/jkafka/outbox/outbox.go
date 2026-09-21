package outbox

import (
	"context"
	"errors"

	"github.com/Freelance-launchpad/backend-interview/common/jkafka"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/outbox/types"
	"github.com/google/uuid"
)

var _ Store = (*store)(nil)
var _ WorkerStore = (*store)(nil)

type Store interface {
	// id is the resource ID (does not have to be unique)
	// state is the state of the resource
	Insert(ctx context.Context, topic string, id string, eventType jkafka.Type, state string) error
}

type WorkerStore interface {
	GetNextEvent(ctx context.Context) (types.StoredEvent, error)
	SetEventSent(ctx context.Context, sentEventID uuid.UUID) error
}

const CronjobCommand = "outbox"

var (
	ErrNoProducerFound = errors.New("no producer found for topic")
	ErrNoEventFound    = errors.New("no event found to send")
)
