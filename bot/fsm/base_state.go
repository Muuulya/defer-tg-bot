package fsm

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/buttons"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/messages"
	"github.com/muuulya/defer-tg-bot/storage"
)

type BaseState struct {
	*AbstractState
}

func NewBaseState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *BaseState {
	return &BaseState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *BaseState) Name() string { return baseStateName }

func (s *BaseState) Enter(user *data.User) error {
	s.showStateMessage(user, "")
	return nil
}

func (s *BaseState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
	if update.Message != nil {
		message := update.Message
		defer s.manager.RemoveMessage(user.ID, message.MessageID)

		if s.isStartCommand(message) {
			return s.handleStartCommand(user, message)
		}

		s.showStateMessage(user, messages.Unknown)
		return "", nil
	}

	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		if s.isButtonPressed(callbackQuery, buttons.Channels) {
			s.manager.SendCallbackMessage(callbackQuery, "Открываем каналы")
			return showMyChannelsStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.AddPost) {
			s.manager.SendCallbackMessage(callbackQuery, "Будем добавлять новый пост")
			return addChannelStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.Settings) {
			s.manager.SendCallbackMessage(callbackQuery, "Этого меню еще нет")
			return "", nil
		}
		if s.isButtonPressed(callbackQuery, buttons.Info) {
			s.manager.SendMessage(user.ID, messages.Info)
			s.manager.RemoveMessage(user.ID, user.CurrentDialogMessageID)
			s.showStateMessage(user, "")
			s.manager.SendCallbackMessage(callbackQuery, "Как все устроено")
			return "", nil
		}
	}

	return "", nil
}

func (s *BaseState) Exit(user *data.User) error {
	return nil
}

func (s *BaseState) showStateMessage(user *data.User, extraMessageText string) {
	text := extraMessageText + messages.EnterBaseState
	keyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			s.getInlineButton(buttons.AddPost),
		),
		tgbotapi.NewInlineKeyboardRow(
			s.getInlineButton(buttons.Channels),
		),
		tgbotapi.NewInlineKeyboardRow(
			s.getInlineButton(buttons.Settings),
			s.getInlineButton(buttons.Info),
		),
	)

	s.editOrCreateNewMessageWithButtons(user, text, keyboard)
}
