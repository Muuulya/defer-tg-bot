package storage

import (
	"time"

	"github.com/muuulya/defer-tg-bot/bot/data"
)

type Storage interface {
	Close()

	AddUser(user *data.User) error
	UpdateUserState(user *data.User) error
	UpdateUserCurrentDialogMessage(user *data.User) error
	UpdateUserSelectedChannelID(user *data.User) error
	UpdateUserCurrentChannelPage(user *data.User) error
	TryGetUser(userID int64) (user *data.User, found bool, error error)

	AddChannel(userID int64, channel *data.Channel) error
	UpdateChannelName(userID int64, channel *data.Channel) error
	TryGetChannel(userID int64, channelID int64) (channel *data.Channel, found bool, error error)
	GetAllChannels(userID int64) ([]data.Channel, error)
	RemoveChannel(userID int64, channelID int64) error

	AddDefferedMessage(defferedMessage *data.DefferedMessage) error
	AddDefferedMessages(defferedMessage []data.DefferedMessage) error
	GetAllDefferedMessages() (defferedMessages []data.DefferedMessage, error error)
	GetNextDefferedMessages(time *time.Time) (defferedMessages []data.DefferedMessage, error error)
	RemoveDefferedMessages(defferedMessages []data.DefferedMessage) error
}
