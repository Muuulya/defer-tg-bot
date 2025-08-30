package bot

import (
	"context"
	"log"
	"time"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/fsm"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/storage"
)

type Bot struct {
	tgAPI   *tgbotapi.BotAPI
	sotrage storage.Storage
	fsm     fsm.StateMashine
}

func NewBot(env *data.ENV) *Bot {
	storage, err := storage.NewStorageDB()
	if err != nil {
		log.Fatal(err)
	}

	tgAPI, err := tgbotapi.NewBotAPI(env.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	tgAPI.Debug = env.DebugMod
	log.Printf("Authorized on account %s", tgAPI.Self.UserName)

	sender, err := manager.NewMessageManager(tgAPI)
	if err != nil {
		log.Fatal(err)
	}

	fsm, err := fsm.NewStateMashine(tgAPI, storage, sender)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{tgAPI: tgAPI, sotrage: storage, fsm: *fsm}
}

func (b *Bot) Start() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := b.tgAPI.GetUpdatesChan(updateConfig)

	user, _, _ := b.sotrage.TryGetUser(201921942)
	SchedulePerChannel(
		context.Background(),
		b.tgAPI,
		b.sotrage,
		user.ID,
		"Тестовая рассылка через 2 минуты",
		2*time.Minute,
		log.Printf,
	)
	// при необходимости: cancel()

	for update := range updates {
		user, found, err := b.getUser(&update)
		if err != nil {
			log.Println(err)
		} else if found {
			b.fsm.Handle(user, &update)
		} else {
			b.createNewUserAndStart(&update)
		}
	}
}

func (b *Bot) Stop() { b.sotrage.Close() }

func (b *Bot) createNewUserAndStart(update *tgbotapi.Update) {
	userID, found := b.getUserID(update)
	if !found {
		return
	}

	userName, found := b.getUserName(update)
	if !found {
		return
	}

	user := data.NewUser(userID, userName)

	err := b.sotrage.AddUser(user)
	if err != nil {
		log.Printf("faild create new user with error: %s", err)
		return
	}

	err = b.fsm.SetStartState(user)
	if err != nil {
		log.Printf("faild create new user with error: %s", err)
		return
	}

	b.fsm.Handle(user, update)
}

func (b *Bot) getUser(update *tgbotapi.Update) (user *data.User, found bool, err error) {
	userID, found := b.getUserID(update)
	if !found {
		return
	}

	user, found, err = b.sotrage.TryGetUser(userID)
	if err != nil {
		log.Printf("faild get user with id: %d with error: %s", userID, err)
		return
	}
	if !found {
		log.Printf("user with id: %d don't found", userID)
		return
	}

	return
}

func (b *Bot) getUserID(update *tgbotapi.Update) (userID int64, found bool) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		return update.Message.Chat.ID, true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat.IsPrivate() {
		return update.CallbackQuery.Message.Chat.ID, true
	}

	return 0, false
}

func (b *Bot) getUserName(update *tgbotapi.Update) (userName string, found bool) {
	if update.Message != nil && update.Message.Chat.IsPrivate() {
		return update.Message.Chat.UserName, true
	}
	if update.CallbackQuery != nil && update.CallbackQuery.Message != nil && update.CallbackQuery.Message.Chat.IsPrivate() {
		return update.CallbackQuery.Message.Chat.UserName, true
	}

	return "", false
}
