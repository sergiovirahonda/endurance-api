package notification

import (
	"github.com/sergiovirahonda/endurance-api/internal/config"
	"gopkg.in/telebot.v3"
)

type TelegramClient struct {
	Token  string
	ChatID int64
	Bot    *telebot.Bot
}

func NewTelegramClient(
	cfg *config.Config,
	token string,
	chatID int64,
) *TelegramClient {
	bot, err := telebot.NewBot(
		telebot.Settings{
			Token: token,
		},
	)
	if err != nil {
		panic(err)
	}
	return &TelegramClient{
		Token:  token,
		Bot:    bot,
		ChatID: chatID,
	}
}
