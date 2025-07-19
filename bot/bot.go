package bot

import (
	"database/sql"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	API *tgbotapi.BotAPI
	DB  *sql.DB
}

func NewBot(api *tgbotapi.BotAPI, db *sql.DB) *Bot {
	return &Bot{API: api, DB: db}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.API.GetUpdatesChan(u)

	// Start scheduler goroutine
	go b.StartScheduler()

	for update := range updates {
		if update.Message != nil {
			b.HandleMessage(update.Message)
		} else if update.CallbackQuery != nil {
			b.HandleCallback(update.CallbackQuery)
		}
	}
}
