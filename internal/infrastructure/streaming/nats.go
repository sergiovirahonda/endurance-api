package streaming

import (
	"fmt"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/sergiovirahonda/endurance-api/internal/config"
)

func NewConnection(cfg *config.Config) *nats.Conn {
	logger := config.GetLogger()
	logger.Info("Creating new NATS connection...")
	nc, err := nats.Connect(
		fmt.Sprintf(
			"nats://%s:%s@%s:%s",
			cfg.NatsUsername,
			cfg.NatsPassword,
			cfg.NatsClusterAddress,
			cfg.NatsClusterPort,
		),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(10),
		nats.ReconnectWait(time.Second),
		nats.Name(cfg.NatsClientID),
	)
	if err != nil {
		panic(fmt.Sprintf("Error connecting to nats: %s", err))
	}
	logger.Info("NATS connection created successfully.")
	return nc
}

func NewStream(nc *nats.Conn, cfg *config.Config) nats.JetStreamContext {
	logger := config.GetLogger()
	// Create JetStream Context
	js, err := nc.JetStream()
	if err != nil {
		logger.Fatal(fmt.Sprintf("Error creating JetStream context: %s", err))
	}
	// Create a stream
	logger.Info("Creating stream...")
	si, err := js.AddStream(
		&nats.StreamConfig{
			Name:        cfg.NATS.NatsStreamName,
			Description: "Nutriflex events stream",
			Subjects: []string{
				cfg.NATS.NatsDomainEventsPattern,
				cfg.NATS.NatsDLQSubject,
			},
			Retention:  nats.WorkQueuePolicy,
			Duplicates: 2 * time.Minute,
		},
	)
	if err != nil {
		logger.Warnf("Error creating stream: %s", err)
	}
	logger.Infof("Stream created: %+v", si)
	return js
}

func NewSubscriber(
	ctx echo.Context,
	js nats.JetStreamContext,
	cfg *config.Config,
	processingFunc func(ctx echo.Context, msg *nats.Msg),
) (*nats.Subscription, error) {
	sub, err := js.QueueSubscribe(
		cfg.NATS.NatsDomainEventsPattern,
		cfg.NATS.NatsClientID,
		func(msg *nats.Msg) {
			processingFunc(ctx, msg)
		},
		nats.Durable(cfg.NATS.NatsClientID),
		nats.AckExplicit(),
		nats.ManualAck(),
	)
	if err != nil {
		log.Fatalf("Failed to subscribe: %v", err)
	}
	return sub, err
}
