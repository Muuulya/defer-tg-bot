package bot

import (
	"context"
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/fsm"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/bot/receiver"
	"github.com/muuulya/defer-tg-bot/bot/sheduler"
	"github.com/muuulya/defer-tg-bot/storage"
	"github.com/muuulya/defer-tg-bot/storage/storageDB"
)

type Bot struct {
	tgAPI    *tgbotapi.BotAPI
	sotrage  storage.Storage
	fsm      *fsm.StateMashine
	receiver *receiver.Receiver
	sheduler *sheduler.Sheduler
}

func NewBot(env *data.ENV) *Bot {
	storage, err := storageDB.NewStorageDB()
	if err != nil {
		log.Fatal(err)
	}

	tgAPI, err := tgbotapi.NewBotAPI(env.BotToken)
	if err != nil {
		log.Fatal(err)
	}
	tgAPI.Debug = env.DebugMod
	log.Printf("Authorized on account %s", tgAPI.Self.UserName)

	manager, err := manager.NewMessageManager(tgAPI)
	if err != nil {
		log.Fatal(err)
	}

	fsm, err := fsm.NewStateMashine(tgAPI, storage, manager)
	if err != nil {
		log.Fatal(err)
	}

	receiver, err := receiver.NewReceiver(tgAPI, storage, fsm)
	if err != nil {
		log.Fatal(err)
	}

	sheduler, err := sheduler.NewSheduler(storage, manager)
	if err != nil {
		log.Fatal(err)
	}

	return &Bot{
		tgAPI:    tgAPI,
		sotrage:  storage,
		fsm:      fsm,
		receiver: receiver,
		sheduler: sheduler,
	}
}

func (b *Bot) Start(ctx context.Context) {
	receiverContext, _ := context.WithCancel(ctx)
	go b.receiver.Start(receiverContext)

	shedulerContext, _ := context.WithCancel(ctx)
	go b.sheduler.Start(shedulerContext)
}

func (b *Bot) Stop() { b.sotrage.Close() }
