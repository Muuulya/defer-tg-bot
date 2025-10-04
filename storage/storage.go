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
	GetUser(userID int64) (user *data.User, err error)
	GetAllUsers() (users []data.User, err error)

	AddChannel(userID int64, channel *data.Channel) error
	UpdateChannelName(userID int64, channel *data.Channel) error
	GetChannel(userID int64, channelID int64) (channel *data.Channel, err error)
	GetAllUserChannels(userID int64) (channels []data.Channel, err error)
	RemoveChannel(userID int64, channelID int64) error

	AddMessage(message *data.DefferedMessage) error
	AddMessagePack(pack *data.DefferedMessagePack) error
	GetMessagePackForUserChannelAfter(userID int64, channelID int64, after time.Time) (pack *data.DefferedMessagePack, err error)
	GetMissedMessagesPacksBefor(befor time.Time) (packs []data.DefferedMessagePack, err error)
	UpdateMessageStatus(message *data.DefferedMessage) error
	UpdateMessagePackStatus(pack *data.DefferedMessagePack) error
	RemoveMessage(message *data.DefferedMessage) error
	RemoveMessagePack(pack *data.DefferedMessagePack) error
}
