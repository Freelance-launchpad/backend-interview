package outbox

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Freelance-launchpad/backend-interview/common/jerror/v2"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka"
	"github.com/Freelance-launchpad/backend-interview/common/jkafka/outbox/types"
	transactor "github.com/Thiht/transactor/stdlib"
	"github.com/google/uuid"
)

// deploy:
// CREATE TYPE event_type AS ENUM ('created', 'updated', 'deleted');
//
// CREATE TABLE IF NOT EXISTS outbox (
//   event_id    UUID       NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
//   topic       TEXT       NOT NULL,
// 	 id          TEXT       NOT NULL,
//   type        event_type NOT NULL,
//   state       TEXT       NOT NULL,
//   sent_at     TIMESTAMP,
//   created_at  TIMESTAMP NOT NULL DEFAULT now()
// );
//
// revert:
// DROP TABLE IF EXISTS outbox;
//
// DROP TYPE event_type;
//
// verify:
// SELECT event_id, topic, id, type, state, sent_at, created_at FROM outbox LIMIT 1;

type store struct {
	dbGetter transactor.DBGetter
}

func NewStore(dbGetter transactor.DBGetter) *store {
	return &store{
		dbGetter: dbGetter,
	}
}

const queryInsert = `
INSERT INTO outbox (topic, id, type, state)
VALUES ($1, $2, $3, $4)`

func (s *store) Insert(ctx context.Context, topic string, id string, eventType jkafka.Type, state string) (err error) {
	defer jerror.Wrap(&err)

	if _, err := s.dbGetter(ctx).ExecContext(ctx, queryInsert, topic, id, eventType, state); err != nil {
		return err
	}

	return nil
}

const queryGetNextEvent = `
SELECT event_id, topic, id, type, state
FROM outbox
WHERE sent_at IS NULL
ORDER BY created_at
LIMIT 1
FOR UPDATE SKIP LOCKED`

func (s *store) GetNextEvent(ctx context.Context) (e types.StoredEvent, err error) {
	defer jerror.Wrap(&err)

	if err := s.dbGetter(ctx).QueryRowContext(ctx, queryGetNextEvent).
		Scan(&e.EventID, &e.Topic, &e.ID, &e.EventType, &e.State); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return types.StoredEvent{}, ErrNoEventFound
		}
		return types.StoredEvent{}, err
	}

	return e, nil
}

const querySetEventSent = `
UPDATE outbox
SET sent_at = now()
WHERE event_id = $1`

func (s *store) SetEventSent(ctx context.Context, sentEventID uuid.UUID) (err error) {
	defer jerror.Wrap(&err)

	if _, err := s.dbGetter(ctx).ExecContext(ctx, querySetEventSent, sentEventID); err != nil {
		return err
	}

	return nil
}
