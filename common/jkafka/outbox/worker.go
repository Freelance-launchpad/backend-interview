package outbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Freelance-launchpad/backend-interview/common/jclock"
	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka"
	"github.com/Thiht/transactor"
)

type Worker struct {
	store      WorkerStore
	transactor transactor.Transactor
	producers  map[string]jkafka.Producer
	clock      jclock.Clock
}

func NewWorker(store WorkerStore, transactor transactor.Transactor, producers map[string]jkafka.Producer) *Worker {
	w := &Worker{
		store:      store,
		transactor: transactor,
		producers:  producers,
		clock:      jclock.RealClock{},
	}
	return w
}

type message struct {
	jkafka.Event
	State json.RawMessage `json:"state"`
}

func (w *Worker) Run(ctx context.Context) (err error) {
	defer jerror.Wrap(&err)

	for {
		if err := w.transactor.WithinTransaction(ctx, func(ctx context.Context) error {
			event, err := w.store.GetNextEvent(ctx)
			if err != nil {
				return err
			}

			producer, ok := w.producers[event.Topic]
			if !ok {
				return fmt.Errorf("topic %s: %w", event.Topic, ErrNoProducerFound)
			}

			sentAt := w.clock.Now()
			if err := producer.Write(ctx, message{
				ID:      event.ID,
				EventAt: sentAt,
				Type:    event.EventType,
				State:   json.RawMessage(event.State),
			}); err != nil {
				return fmt.Errorf("failed to send event %s: %w", event.EventID, err)
			}

			return w.store.SetEventSent(ctx, event.EventID)
		}); err != nil {
			if errors.Is(err, ErrNoEventFound) {
				return nil
			}
			return err
		}
	}
}
