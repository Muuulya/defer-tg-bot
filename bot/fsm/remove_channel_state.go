package fsm

import (
	"fmt"
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/buttons"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/messages"
	"github.com/muuulya/defer-tg-bot/storage"
)

type RemoveChannelState struct {
	*AbstractState
}

func NewRemoveChannelState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *RemoveChannelState {
	return &RemoveChannelState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *RemoveChannelState) Name() string { return removeChannelStateName }

func (s *RemoveChannelState) Enter(user *data.User) error {
	return s.showStateMessage(user, "")
}

func (s RemoveChannelState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
	if update.Message != nil {
		message := update.Message
		defer s.manager.RemoveMessage(user.ID(), message.MessageID)

		if s.isStartCommand(message) {
			return s.handleStartCommand(user, message)
		}

		err = s.showStateMessage(user, messages.Unknown)
		return "", err
	}

	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		if s.isButtonPressed(callbackQuery, buttons.Cancel) {
			s.manager.SendCallbackMessage(callbackQuery, "Отмена")
			return showMyChannelsStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.RemoveChannel) {
			if err = s.storage.RemoveChannel(user.ID(), user.SelectedChannelID()); err != nil {
				log.Println(err)
				s.manager.SendCallbackMessage(callbackQuery, "Что-то пошло не так. Не удалось удалить канал")
			} else {
				s.manager.SendCallbackMessage(callbackQuery, "Канал удален")
			}
			return showMyChannelsStateName, nil
		}
	}

	return "", nil
}

func (s *RemoveChannelState) Exit(user *data.User) error {
	user.SetSelectedChannel(0)
	s.storage.UpdateUserSelectedChannelID(user)
	return nil
}

func (s *RemoveChannelState) showStateMessage(user *data.User, extraMessageText string) error {
	text := extraMessageText
	var keyboard tgbotapi.InlineKeyboardMarkup

	channel, err := s.storage.GetChannel(user.ID(), user.SelectedChannelID())
	if err != nil {
		text += messages.ChannelNotFound
		keyboard = s.getInlineKeyboard(
			[]buttons.Button{buttons.Cancel},
		)
	} else {
		text += fmt.Sprintf(messages.RemoveChannel, channel.Name)
		keyboard = s.getInlineKeyboard(
			[]buttons.Button{buttons.RemoveChannel},
			[]buttons.Button{buttons.Cancel},
		)
	}

	s.editOrCreateNewMessageWithButtons(user, text, keyboard)
	return err
}
