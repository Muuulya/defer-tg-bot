package fsm

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"github.com/muuulya/defer-tg-bot/bot/buttons"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/storage"
)

type State interface {
	Name() string
	Enter(user *data.User) error
	Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error)
	Exit(user *data.User) error
}

type AbstractState struct {
	tgAPI   *tgbotapi.BotAPI
	manager *manager.MessageManager
	storage storage.Storage
}

func NewAbstractState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *AbstractState {
	return &AbstractState{tgAPI: tgAPI, manager: manager, storage: storage}
}

func (s *AbstractState) handleStartCommand(_ *data.User, _ *tgbotapi.Message) (nextState string, err error) {
	return baseStateName, nil
}

func (s *AbstractState) isStartCommand(message *tgbotapi.Message) bool {
	return message.IsCommand() && message.Command() == "start"
}

func (s *AbstractState) editOrCreateNewMessage(user *data.User, text string) {
	_, err := s.manager.EditMessage(user.ID, user.CurrentDialogMessageID, text)
	if err != nil {
		newm, err := s.manager.SendMessage(user.ID, text)
		if err == nil {
			s.manager.RemoveMessage(user.ID, user.CurrentDialogMessageID)
			user.CurrentDialogMessageID = newm.MessageID
			s.storage.UpdateUserCurrentDialogMessage(user)
		}
	}
}

func (s *AbstractState) editOrCreateNewMessageWithButtons(user *data.User, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	_, err := s.manager.EditMessageWithInlineButtons(user.ID, user.CurrentDialogMessageID, text, keyboard)
	if err != nil {
		newm, err := s.manager.SendMessageWithInlineButtons(user.ID, text, keyboard)
		if err == nil {
			s.manager.RemoveMessage(user.ID, user.CurrentDialogMessageID)
			user.CurrentDialogMessageID = newm.MessageID
			s.storage.UpdateUserCurrentDialogMessage(user)
		}
	}
}

func (s *AbstractState) getInlineButton(button buttons.Button) tgbotapi.InlineKeyboardButton {
	return tgbotapi.NewInlineKeyboardButtonData(button.Name, button.Value)
}

func (s *AbstractState) getInlineButtonRow(buttons ...buttons.Button) []tgbotapi.InlineKeyboardButton {
	rows := make([]tgbotapi.InlineKeyboardButton, 0, len(buttons))
	for _, b := range buttons {
		rows = append(rows, s.getInlineButton(b))
	}
	return tgbotapi.NewInlineKeyboardRow(rows...)
}

func (s *AbstractState) getInlineKeyboard(rows ...[]buttons.Button) tgbotapi.InlineKeyboardMarkup {
	keyboard := make([][]tgbotapi.InlineKeyboardButton, 0, len(rows))
	for _, row := range rows {
		keyboard = append(keyboard, s.getInlineButtonRow(row...))
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboard...)
}

func (s *AbstractState) isButtonPressed(callback *tgbotapi.CallbackQuery, button buttons.Button) bool {
	return callback.Data == button.Value
}
