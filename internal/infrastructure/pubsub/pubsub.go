package pubsub

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/sergiovirahonda/endurance-api/internal/config"
)

type EventsPubSub struct {
	NatsConnection  *nats.Conn
	JetStreamStream nats.JetStreamContext
	Subscriber      *nats.Subscription
}

// Factories

func NewEventsPubSub(
	nc *nats.Conn,
	js nats.JetStreamContext,
) *EventsPubSub {
	return &EventsPubSub{
		NatsConnection:  nc,
		JetStreamStream: js,
	}
}

// Producer receivers

func (ps *EventsPubSub) Publish(subject string, data interface{}) error {
	event, err := json.Marshal(data)
	if err != nil {
		return err
	}
	eventID := uuid.New().String()
	_, err = ps.JetStreamStream.PublishAsync(subject, event, nats.MsgId(eventID))
	if err != nil {
		return err
	}
	return nil
}

func (ps *EventsPubSub) SendToDLQ(data interface{}) error {
	conf := config.GetConfig()
	event, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = ps.JetStreamStream.PublishAsync(conf.NatsDLQSubject, event)
	if err != nil {
		return err
	}
	return nil
}

// Consumer receivers

func (ps *EventsPubSub) FetchFromSubjects() (*nats.Msg, error) {
	msgs, err := ps.Subscriber.Fetch(1, nats.MaxWait(5*time.Second))
	if err != nil {
		return nil, err
	}
	return msgs[0], nil
}
