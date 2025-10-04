package handlers

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/data"
)

type MessageHandler interface {
	Name() string
	Handle(user *data.User, message *tgbotapi.Message) (nextState string, err error)
}
