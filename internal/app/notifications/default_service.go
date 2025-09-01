package notifications

import (
	"fmt"
	"strconv"

	"github.com/labstack/echo/v4"
	keys "github.com/sergiovirahonda/endurance-api/internal/app/key"
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"github.com/sergiovirahonda/endurance-api/internal/infrastructure/notification"
	"gopkg.in/telebot.v3"
)

// Structs

type DefaultNotificationService struct {
	KeyService     keys.KeyService
	TelegramClient *notification.TelegramClient
}

// Factories

func NewDefaultNotificationService(
	keyService keys.KeyService,
	telegramClient *notification.TelegramClient,
) *DefaultNotificationService {
	return &DefaultNotificationService{
		KeyService:     keyService,
		TelegramClient: telegramClient,
	}
}

// Receivers

func (s *DefaultNotificationService) SendMessage(
	ctx echo.Context,
	message string,
) error {
	keys, err := s.KeyService.GetTelegramKeys(ctx)
	if err != nil {
		return err
	}
	s.TelegramClient.Bot.Token = keys.Key
	secret, err := strconv.ParseInt(keys.Secret, 10, 64)
	if err != nil {
		return err
	}
	msg, err := s.TelegramClient.Bot.Send(
		&telebot.User{ID: secret},
		message,
	)
	logger := config.GetLogger()
	logger.Infof("Telegram message sent: %+v", msg.Text)
	if err != nil {
		return err
	}
	return nil
}

// Integrations

func (s *DefaultNotificationService) SendTradeNotification(
	ctx echo.Context,
	originSymbol string,
	newSymbol string,
	entryPrice float64,
	profit float64,
	profitPercentage float64,
) error {
	message := fmt.Sprintf(
		"❗ Trade operation executed.\n\n"+
			"💰 %s >> %s\n"+
			"- Entry price: %f\n"+
			"- Profit: %f USDT\n"+
			"- Profit percentage: %f\n",
		originSymbol,
		newSymbol,
		entryPrice,
		profit,
		profitPercentage,
	)
	return s.SendMessage(ctx, message)
}

func (s *DefaultNotificationService) SendStopLossNotification(
	ctx echo.Context,
	originSymbol string,
	stopLossPrice float64,
	loss float64,
	lossPercentage float64,
) error {
	message := fmt.Sprintf(
		"❗ Stop loss triggered.\n\n"+
			"💰 %s >> USDT\n"+
			"- Stop loss price: %f\n"+
			"- Loss: %f USDT\n"+
			"- Loss percentage: %f\n",
		originSymbol,
		stopLossPrice,
		loss,
		lossPercentage,
	)
	return s.SendMessage(ctx, message)
}
