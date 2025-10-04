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

type ShowChannelState struct {
	*AbstractState
}

func NewShowChannelState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *ShowChannelState {
	return &ShowChannelState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *ShowChannelState) Name() string { return showChannelStateName }

func (s *ShowChannelState) Enter(user *data.User) error {
	return s.showStateMessage(user, "")
}

func (s *ShowChannelState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
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
		if s.isButtonPressed(callbackQuery, buttons.UpdateChannel) {
			if err = s.updateChannelName(user); err != nil {
				log.Println(err)
				s.manager.SendCallbackMessage(callbackQuery, "Не удалось обновить имя")
			} else {
				s.manager.SendCallbackMessage(callbackQuery, "Имя успешно обновлено")
			}
			err = s.showStateMessage(user, "")
			return "", err
		}
		if s.isButtonPressed(callbackQuery, buttons.RemoveChannel) {
			s.manager.SendCallbackMessage(callbackQuery, "Ты уверен, что хочушь удалить канал?")
			return removeChannelStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.Return) {
			s.manager.SendCallbackMessage(callbackQuery, "Летим назад")
			return showMyChannelsStateName, nil
		}
	}

	return "", nil
}

func (s *ShowChannelState) Exit(user *data.User) error {
	return nil
}

func (s *ShowChannelState) showStateMessage(user *data.User, extraMessageText string) error {
	text := extraMessageText
	var keyboard tgbotapi.InlineKeyboardMarkup

	channel, err := s.storage.GetChannel(user.ID(), user.SelectedChannelID())
	if err != nil {
		text = messages.ChannelNotFound
		keyboard = s.getInlineKeyboard(
			[]buttons.Button{buttons.Return},
		)
	} else {
		text += fmt.Sprintf(messages.RemoveChannel, channel.Name)
		keyboard = s.getInlineKeyboard(
			[]buttons.Button{buttons.UpdateChannel},
			[]buttons.Button{buttons.RemoveChannel},
			[]buttons.Button{buttons.Return},
		)
	}

	s.editOrCreateNewMessageWithButtons(user, text, keyboard)
	return err
}

func (s *ShowChannelState) updateChannelName(user *data.User) error {
	chatInfoConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: user.SelectedChannelID()}}
	if chat, err := s.tgAPI.GetChat(chatInfoConfig); err != nil {
		return err
	} else {
		channel := data.NewChannel(chat.ID, chat.Title)
		if err = s.storage.UpdateChannelName(user.ID(), channel); err != nil {
			return err
		}
	}

	return nil
}
