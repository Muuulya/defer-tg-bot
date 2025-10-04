package sheduler

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/bot/manager"
	"github.com/muuulya/defer-tg-bot/storage"
)

type Sheduler struct {
	storage            storage.Storage
	manager            *manager.MessageManager
	shedulChannelPacks map[data.UserChannelPair]context.CancelFunc
	mutex              sync.Mutex
}

func NewSheduler(storage storage.Storage, manager *manager.MessageManager) (sheduler *Sheduler, err error) {
	sheduler = &Sheduler{
		storage:            storage,
		manager:            manager,
		shedulChannelPacks: make(map[data.UserChannelPair]context.CancelFunc),
		mutex:              sync.Mutex{},
	}

	return sheduler, nil
}

func (s *Sheduler) Start(ctx context.Context) {
	users, err := s.storage.GetAllUsers()
	if err != nil {
		log.Println(err)
	}

	s.checkMissedMessages(ctx, users)
	s.scheduleMessages(ctx, users)
}

func (s *Sheduler) checkMissedMessages(ctx context.Context, users []data.User) {
	packs, err := s.storage.GetMissedMessagesPacksBefor(s.now())
	if err != nil {
		log.Println(err)
	}

	for _, pack := range packs {
		if err = s.sendMessagePack(pack); err != nil {
			log.Println(err)
		}
	}
}

func (s *Sheduler) scheduleMessages(ctx context.Context, users []data.User) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, user := range users {
		channels, err := s.storage.GetAllUserChannels(user.ID())
		if err != nil {
			log.Println(err)
		}

		for _, channel := range channels {
			pack, err := s.storage.GetMessagePackForUserChannelAfter(user.ID(), channel.ID(), s.now())
			if err != nil {
				log.Println(err)
			}

			pair := data.NewUserChannelPair(user.ID(), channel.ID())
			packContext, cancel := context.WithCancel(ctx)
			go s.scheduleMessagePackForUserChannel(packContext, *pack, pair)
			s.shedulChannelPacks[pair] = cancel
		}
	}
}

func (s *Sheduler) scheduleMessagePackForUserChannel(ctx context.Context, pack data.DefferedMessagePack, pair data.UserChannelPair) {
	timer := s.timerTo(pack.PostedTime())

	for {
		select {
		case <-ctx.Done():
			log.Println("Cancel shedule message pack") // todo решить что нужно или не нужно тут делать.
			return
		case <-timer.C:
			s.mutex.Lock()
			defer s.mutex.Unlock()
			s.shedulChannelPacks[pair] = nil
			if err := s.sendMessagePack(pack); err != nil {
				log.Println(err)
			}
		}
	}
}

func (s *Sheduler) sendMessagePack(pack data.DefferedMessagePack) error {
	var counter = 0
	for i, message := range pack.Messages() {
		if _, err := s.manager.CopyMessage(message.UserID(), message.MessageID(), message.TargetChannelID()); err != nil {
			counter++
			log.Println(err)
		}
		pack.SetMessageStatus(i, true)
	}

	if err := s.storage.UpdateMessagePackStatus(&pack); err != nil {
		log.Println(err)
	}

	if counter > 0 {
		return ErrorSomeMessageNotSend //todo подумать над тем что возвращать и как обрабатывать ошибки.
	}

	return nil
}

func (s *Sheduler) now() time.Time {
	return time.Now().UTC()
}

func (s *Sheduler) duration(t time.Time) time.Duration {
	return t.Sub(s.now())
}

func (s *Sheduler) timerTo(t time.Time) *time.Timer {
	return time.NewTimer(s.duration(t))
}
