package events

import (
	"time"

	"github.com/google/uuid"
)

type BaseEvent struct {
	ID        uuid.UUID `json:"id"`
	Domain    string    `json:"domain"`
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
}

// Factories

type BaseEventFactory struct{}

// Factory receivers

func (f BaseEventFactory) NewBaseEvent(domain, eventType string) *BaseEvent {
	return &BaseEvent{
		ID:        uuid.New(),
		Domain:    domain,
		Type:      eventType,
		Timestamp: time.Now().UTC(),
	}
}
