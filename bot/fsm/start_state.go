package fsm

import (
	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/messages"
	"github.com/muuulya/defer-tg-bot/storage"
)

type StartState struct {
	*AbstractState
}

func NewStartState(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) *StartState {
	return &StartState{AbstractState: NewAbstractState(tgAPI, manager, storage)}
}

func (s *StartState) Name() string { return startStateName }

func (s *StartState) Enter(user *data.User) error {
	_, err := s.manager.SendMessage(user.ID, messages.Hello)
	if err != nil {
		return err
	}
	_, err = s.manager.SendMessage(user.ID, messages.Info)
	if err != nil {
		return err
	}
	return nil
}

func (s *StartState) Handle(user *data.User, update *tgbotapi.Update) (nextState string, err error) {
	return baseStateName, nil
}

func (s *StartState) Exit(user *data.User) error { return nil }
