package messaging

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/nats.go"
	"github.com/sergiovirahonda/endurance-api/internal/app/markets"
	uacs "github.com/sergiovirahonda/endurance-api/internal/app/uac"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/pubsub"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/streaming"
)

type DefaultMessagingService struct {
	eventsPubSub            pubsub.EventsPubSub
	uacService              uacs.UacService
	marketDataEventRegistry markets.MarketDataEventRegistry
}

func NewDefaultMessagingService(
	eventsPubSub pubsub.EventsPubSub,
	uacService uacs.UacService,
	marketDataEventRegistry markets.MarketDataEventRegistry,
) *DefaultMessagingService {
	return &DefaultMessagingService{
		eventsPubSub:            eventsPubSub,
		uacService:              uacService,
		marketDataEventRegistry: marketDataEventRegistry,
	}
}

func (ms *DefaultMessagingService) EventsLoop() {
	logger := config.GetLogger()
	conf := config.GetConfig()
	ctx := echo.New().NewContext(nil, nil)
	logger.Info("Starting events loop...")
	ctx.Set("logger", logger)
	sub, err := streaming.NewSubscriber(
		ctx,
		ms.eventsPubSub.JetStreamStream,
		conf,
		ms.ProcessMessage,
	)
	if err != nil {
		logger.Fatalf("Error creating subscriber: %s", err)
	}
	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	logger.Info("Shutting down subscriber...")
	err = sub.Unsubscribe()
	if err != nil {
		logger.Fatalf(fmt.Sprintf("Error unsubscribing: %s", err))
	}
	logger.Info("Subscriber shut down.")
}

func (ms DefaultMessagingService) ProcessMessage(ctx echo.Context, msg *nats.Msg) {
	conf := config.GetConfig()
	logger := config.GetLogger()
	logger.Info("Received event from subject: ", msg.Subject)
	DomainEventSubjectRoot := strings.Replace(
		conf.NATS.NatsDomainEventsPattern,
		"*",
		"",
		-1,
	)
	if !strings.Contains(msg.Subject, DomainEventSubjectRoot) {
		logger.Warn("Unknown subject: ", msg.Subject)
		err := msg.Nak()
		if err != nil {
			logger.Errorf(fmt.Sprintf("Error sending NAK: %s", err))
		}
	} else {
		err := ms.HandleDomainEvent(ctx, msg.Data)
		if err != nil {
			logger.Error(fmt.Sprintf("Error handling domain event (%s): %s", msg.Data, err))
			err = ms.eventsPubSub.SendToDLQ(msg.Data)
			if err != nil {
				logger.Error(fmt.Sprintf("Error sending event to DLQ: %s", err))
			}
		}
		err = msg.Ack()
		if err != nil {
			logger.Errorf(fmt.Sprintf("Error sending ACK: %s", err))
		}
	}
}

// Higher-level integration receivers (point of integration between all services)

func (ms DefaultMessagingService) HandleDomainEvent(ctx echo.Context, msg []byte) error {
	// Derive domain event to subscription event handling
	err := ms.marketDataEventRegistry.HandleEvent(ctx, msg)
	if err != nil {
		return err
	}
	// NOTE: Handler other events below
	return nil
}
