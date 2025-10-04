package fsm

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/buttons"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/messages"
	"github.com/muuulya/defer-tg-bot/storage"
)

var channelReStrict = regexp.MustCompile(
	`(?i)(?:https?://)?(?:t\.me|telegram\.me)/(?:s/)?@?([a-z0-9_]{5,32})|(?:^|[^a-z0-9_])@([a-z0-9_]{5,32})`)

type AddChannelState struct {
	*AbstractState
}

func NewAddChannelState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *AddChannelState {
	return &AddChannelState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *AddChannelState) Name() string { return addChannelStateName }

func (s *AddChannelState) Enter(user *data.User) error {
	return s.showStateMessage(user, "")
}

func (s *AddChannelState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
	if update.Message != nil {
		message := update.Message
		defer s.manager.RemoveMessage(user.ID(), message.MessageID)

		if s.isStartCommand(message) {
			return s.handleStartCommand(user, message)
		}

		if channel, err := s.ResolveChannel(message); err == nil {
			getChatMemberConfig := tgbotapi.NewGetChatMember(channel.ID(), s.tgAPI.Self.ID)
			if member, err := s.tgAPI.GetChatMember(getChatMemberConfig); err == nil {
				if member.Status != "administrator" && member.Status != "creator" {
					err = s.showStateMessage(user, messages.BotNotAdmin)
					return "", err
				}
				if err = s.storage.AddChannel(user.ID(), channel); err == nil {
					return showMyChannelsStateName, nil
				} else {
					log.Println(err)
					err = s.showStateMessage(user, messages.AddChannelError)
					return "", err
				}
			} else {
				log.Println(err)
				err = s.showStateMessage(user, messages.BotNotChannelMember)
				return "", err
			}
		} else {
			log.Println(err)
			err = s.showStateMessage(user, messages.BotNameNotFound)
			return "", err
		}
	}

	if update.CallbackQuery != nil {
		callbackQuery := update.CallbackQuery
		if s.isButtonPressed(callbackQuery, buttons.Cancel) {
			return showMyChannelsStateName, nil
		}
	}

	return "", nil
}

func (s *AddChannelState) Exit(user *data.User) error {
	return nil
}

func (s *AddChannelState) ResolveChannel(message *tgbotapi.Message) (channel *data.Channel, err error) {
	channel = &data.Channel{}

	if message.ForwardOrigin != nil && message.ForwardOrigin.IsChannel() {
		channel = data.NewChannel(message.ForwardOrigin.Chat.ID, message.ForwardOrigin.Chat.Title)
		return channel, nil
	}

	text := strings.TrimSpace(message.Text)
	if text == "" {
		return channel, errors.New("text for parsing is empty")
	}

	if found, channelName := s.tryGetChannelNameFromText(text); found {
		chatInfoConfig := tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChannelUsername: "@" + channelName}}
		if chat, err := s.tgAPI.GetChat(chatInfoConfig); err == nil {
			channel = data.NewChannel(chat.ID, chat.Title)
			return channel, nil
		} else {
			return channel, fmt.Errorf("get chat from username: %w", err)
		}
	} else {
		return channel, errors.New("chat not found. try again")
	}
}

func (s *AddChannelState) tryGetChannelNameFromText(text string) (found bool, channelName string) {

	matches := channelReStrict.FindStringSubmatch(text)
	if len(matches) > 2 {
		if matches[1] != "" {
			return true, matches[1]
		}
		if matches[2] != "" {
			return true, matches[2]
		}
	}

	return false, ""
}

func (s *AddChannelState) showStateMessage(user *data.User, extraMessageText string) error {
	text := extraMessageText

	keyboard := s.getInlineKeyboard(
		[]buttons.Button{buttons.Cancel},
	)

	text += messages.AddChannel
	s.editOrCreateNewMessageWithButtons(user, text, keyboard)

	return nil
}
