package fsm

import (
	"errors"
	"fmt"
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/storage"
)

type StateMashine struct {
	states  map[string]State
	storage storage.Storage
	manager *manager.MessageManager
}

func NewStateMashine(
	tgAPI *tgbotapi.BotAPI,
	storage storage.Storage,
	manager *manager.MessageManager,
) (*StateMashine, error) {
	sm := StateMashine{
		states:  make(map[string]State),
		storage: storage,
		manager: manager,
	}
	sm.registerStates(tgAPI, manager, storage)
	return &sm, nil
}

func (sm *StateMashine) registerStates(
	tgAPI *tgbotapi.BotAPI,
	manager *manager.MessageManager,
	storage storage.Storage,
) {
	sm.registerState(NewStartState(tgAPI, manager, storage))
	sm.registerState(NewBaseState(tgAPI, manager, storage))
	sm.registerState(NewShowMyChannelsState(tgAPI, manager, storage))
	sm.registerState(NewAddChannelState(tgAPI, manager, storage))
	sm.registerState(NewShowChannelState(tgAPI, manager, storage))
	sm.registerState(NewRemoveChannelState(tgAPI, manager, storage))
}

func (sm *StateMashine) registerState(state State) {
	log.Printf("state name: %s", state.Name())
	sm.states[state.Name()] = state
}

func (sm *StateMashine) SetStartState(user *data.User) error {
	newState, ok := sm.states[startStateName]
	if !ok {
		return fmt.Errorf("unknown next state: %s", startStateName)
	}

	user.SetCurrentState(startStateName)
	sm.storage.UpdateUserState(user)

	err := newState.Enter(user)
	if err != nil {
		log.Println(err)
		return errors.New("coldn't change the state")
	}

	return nil
}

func (sm *StateMashine) SwitchState(user *data.User, nextStateName string) error {

	if user.CurrentStateName() == "" {
		return errors.New("current state name empty")
	}

	if nextStateName == "" {
		return errors.New("next state name empty")
	}

	currentState, ok := sm.states[user.CurrentStateName()]
	if !ok {
		return fmt.Errorf("unknown current state: %s", user.CurrentStateName)
	}

	newState, ok := sm.states[nextStateName]
	if !ok {
		return fmt.Errorf("unknown next state: %s", nextStateName)
	}

	err := currentState.Exit(user)
	if err != nil {
		log.Println(err)
		return errors.New("coldn't change the state")
	}

	user.SetCurrentState(nextStateName)
	sm.storage.UpdateUserState(user)

	err = newState.Enter(user)
	if err != nil {
		log.Println(err)
		return errors.New("coldn't change the state")
	}

	return nil
}

func (sm *StateMashine) Handle(user *data.User, update *tgbotapi.Update) error {
	if user.CurrentStateName() == "" {
		return errors.New("current state name empty")
	}

	currentState, ok := sm.states[user.CurrentStateName()]
	if !ok {
		return fmt.Errorf("unknown state: %s", user.CurrentStateName)
	}

	nextState, err := currentState.Handle(user, update)
	if err != nil {
		log.Println(err)
		return errors.New("message don't handled")
	}

	if nextState != "" {
		sm.SwitchState(user, nextState)
	}

	return nil
}
