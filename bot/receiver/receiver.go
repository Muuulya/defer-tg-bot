package receiver

import (
	"context"
	"errors"
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/fsm"
	"github.com/muuulya/defer-tg-bot/storage"
)

type Receiver struct {
	tgAPI   *tgbotapi.BotAPI
	sotrage storage.Storage
	fsm     *fsm.StateMashine
}

func NewReceiver(tgAPI *tgbotapi.BotAPI, storage storage.Storage, fsm *fsm.StateMashine) (receiver *Receiver, err error) {
	receiver = &Receiver{
		tgAPI:   tgAPI,
		sotrage: storage,
		fsm:     fsm,
	}

	return receiver, nil
}

func (r *Receiver) Start(ctx context.Context) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updateContext, _ := context.WithCancel(ctx)
	go r.updade(updateContext, updateConfig)
}

func (r *Receiver) updade(ctx context.Context, updateConfig tgbotapi.UpdateConfig) {
	updates := r.tgAPI.GetUpdatesChan(updateConfig)

	for update := range updates {
		user, err := r.getUser(&update)
		if errors.Is(err, storage.ErrorUserNotFound) {
			r.createNewUserAndStart(&update)
		} else if err != nil {
			log.Println(err)
		} else {
			r.fsm.Handle(user, &update)
		}
	}
}

func (r *Receiver) createNewUserAndStart(update *tgbotapi.Update) {
	userID, err := r.getUserID(update)
	if errors.Is(err, ErrorUserIDNotFound) {
		return
	}

	userName, err := r.getUserName(update)
	if errors.Is(err, ErrorUserNameNotFound) {
		return
	}

	user := data.NewUser(userID, userName)

	err = r.sotrage.AddUser(user)
	if err != nil {
		log.Printf("faild create new user with error: %s", err)
		return
	}

	err = r.fsm.SetStartState(user)
	if err != nil {
		log.Printf("faild create new user with error: %s", err)
		return
	}

	r.fsm.Handle(user, update)
}

func (b *Receiver) getUser(update *tgbotapi.Update) (user *data.User, err error) {
	userID, err := b.getUserID(update)
	if err != nil {
		log.Println(err)
		return user, err
	}

	return b.sotrage.GetUser(userID)
}

func (r *Receiver) getUserID(update *tgbotapi.Update) (userID int64, err error) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		userID = update.Message.Chat.ID
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat.IsPrivate() {
		userID = update.CallbackQuery.Message.Chat.ID
	} else {
		err = ErrorUserIDNotFound
	}

	return userID, err
}

func (r *Receiver) getUserName(update *tgbotapi.Update) (userName string, err error) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		userName = update.Message.Chat.UserName
	} else if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat.IsPrivate() {
		userName = update.CallbackQuery.Message.Chat.UserName
	} else {
		err = ErrorUserNameNotFound
	}

	return userName, err
}
