package fsm

import (
	"fmt"
	"log"
	"strconv"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/buttons"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/messages"
	"github.com/muuulya/defer-tg-bot/storage"
)

type ShowMyChannelsState struct {
	*AbstractState
}

func NewShowMyChannelsState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *ShowMyChannelsState {
	return &ShowMyChannelsState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *ShowMyChannelsState) Name() string { return showMyChannelsStateName }

func (s *ShowMyChannelsState) Enter(user *data.User) error {
	s.showStateMessage(user, "")
	return nil
}

func (s *ShowMyChannelsState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
	if update.Message != nil {
		message := update.Message
		defer s.manager.RemoveMessage(user.ID(), message.MessageID)

		if s.isStartCommand(message) {
			return s.handleStartCommand(user, message)
		}

		s.showStateMessage(user, messages.Unknown)
		return "", nil
	}

	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		if s.isButtonPressed(callbackQuery, buttons.Return) {
			s.manager.SendCallbackMessage(callbackQuery, "Назад")
			return baseStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.AddChannel) {
			s.manager.SendCallbackMessage(callbackQuery, "Давай добавим новый канал")
			return addChannelStateName, nil
		}
		if s.isButtonPressed(callbackQuery, buttons.Previous) {
			user.SetChannelPage(user.CurrentChannelPage() - 1)
			s.storage.UpdateUserCurrentChannelPage(user)
			s.manager.SendCallbackMessage(callbackQuery, "Предыдущая страница")
			s.showStateMessage(user, "")
			return "", nil
		}
		if s.isButtonPressed(callbackQuery, buttons.Next) {
			user.SetChannelPage(user.CurrentChannelPage() + 1)
			s.storage.UpdateUserCurrentChannelPage(user)
			s.manager.SendCallbackMessage(callbackQuery, "Следующая страница")
			s.showStateMessage(user, "")
			return "", nil
		}

		channels, err := s.storage.GetAllUserChannels(user.ID())

		if err == nil {
			if isPressed, channelID := s.isChannelButtonPresed(callbackQuery, &channels); isPressed {
				user.SetSelectedChannel(channelID)
				err = s.storage.UpdateUserSelectedChannelID(user)
				if err == nil {
					s.manager.SendCallbackMessage(callbackQuery, "Переходим к каналу")
					return showChannelStateName, nil
				} else {
					log.Println(err)
					s.showStateMessage(user, messages.CommonError)
					return "", err
				}
			}
		} else {
			log.Println(err)
			s.showStateMessage(user, messages.CommonError)
			return "", err
		}
	}

	return "", nil
}

func (s *ShowMyChannelsState) Exit(user *data.User) error {
	user.SetChannelPage(0)
	err := s.storage.UpdateUserCurrentChannelPage(user)
	return err
}

func (s *ShowMyChannelsState) showStateMessage(user *data.User, extraMessageText string) {
	text := extraMessageText
	var keyboard tgbotapi.InlineKeyboardMarkup

	channels, err := s.storage.GetAllUserChannels(user.ID())
	if err == nil {
		if len(channels) == 0 {
			text += messages.NoChannels
		} else {
			channelsCount := len(channels)
			pages := (channelsCount + channelsOnPage - 1) / channelsOnPage
			if user.CurrentChannelPage() < 0 || user.CurrentChannelPage() >= pages {
				user.SetChannelPage(0)
				s.storage.UpdateUserCurrentChannelPage(user)
			}
			text += fmt.Sprintf(messages.AllYourChannel, channelsCount)
			if pages > 1 {
				text += fmt.Sprintf(messages.ChannelPage, user.CurrentChannelPage()+1, pages)
			}
			text += messages.SelectChannel
			keyboard = s.getChannelsButtons(user, channels)
		}

	} else {
		text += messages.ChannelsNotFound
	}

	keyboard.InlineKeyboard = append(keyboard.InlineKeyboard,
		s.getInlineKeyboard(
			[]buttons.Button{buttons.AddChannel},
			[]buttons.Button{buttons.Return},
		).InlineKeyboard...,
	)

	s.editOrCreateNewMessageWithButtons(user, text, keyboard)
}

func (s *ShowMyChannelsState) getChannelsButtons(user *data.User, channels []data.Channel) tgbotapi.InlineKeyboardMarkup {
	total := len(channels)
	page := user.CurrentChannelPage()
	start := page * channelsOnPage
	end := start + channelsOnPage
	if end > total {
		end = total
	}

	rows := make([][]tgbotapi.InlineKeyboardButton, 0, channelsOnPage+1)
	for i, ch := range (channels)[start:end] {
		stringID := strconv.FormatInt(ch.ID(), 10)
		row := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d. %s", start+i+1, ch.Name()), stringID),
		)
		rows = append(rows, row)
	}

	nav := []tgbotapi.InlineKeyboardButton{}
	if page > 0 {
		nav = append(nav, s.getInlineButton(buttons.Previous))
	}
	if page < channelsOnPage {
		nav = append(nav, s.getInlineButton(buttons.Next))
	}
	if len(nav) > 0 {
		rows = append(rows, nav)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (s *ShowMyChannelsState) isChannelButtonPresed(
	callbackQuery *tgbotapi.CallbackQuery,
	channels *[]data.Channel,
) (isPresed bool, channelID int64) {
	if id, err := strconv.ParseInt(callbackQuery.Data, 10, 64); err == nil {
		for _, ch := range *channels {
			if id == ch.ID() {
				return true, ch.ID()
			}
		}
	}
	return false, 0
}
