package jkafka

import "time"

type Message interface {
	Key() string
}

var _ Message = Event{}

// Event is the base type for all events sent to Kafka.
// It contains the common fields that all events must have.
type Event struct {
	ID      string    `json:"id"`
	EventAt time.Time `json:"event_at"`
	Type    Type      `json:"type"`
}

func (e Event) Key() string {
	return e.ID
}

type Type string

const (
	TypeCreated Type = "created"
	TypeUpdated Type = "updated"
	TypeDeleted Type = "deleted"
)
