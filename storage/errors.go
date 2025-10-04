package storage

import "errors"

var (
	//Common
	ErrorEmptyDataTime = errors.New("datatime is empty")

	// Storage
	ErrorCreateStorage = errors.New("failed to create storage")
	ErrorInitStorage   = errors.New("failed to init storage")

	// User
	ErrorAddUser      = errors.New("failed to add a user")
	ErrorUpdateUser   = errors.New("failed to update user")
	ErrorUserNotFound = errors.New("user not found")
	ErrorGetUser      = errors.New("failed to get a user")
	ErrorGetAllUsers  = errors.New("failed to get all users")

	// Channel
	ErrorAddChannel      = errors.New("failed to add a channel")
	ErrorUpdateChannel   = errors.New("failed to update a channel")
	ErrorChannelNotFound = errors.New("channel not found")
	ErrorGetChannel      = errors.New("failed to get channel")
	ErrorGetChannels     = errors.New("failed to get channels")
	ErrorRemoveChannel   = errors.New("failed to remove channel")

	// Message
	ErrorAddMessage          = errors.New("failed to add message")
	ErrorMessageNotFound     = errors.New("message not found")
	ErrorGetMessages         = errors.New("failed to get messages")
	ErrorSeveralMessagePacks = errors.New("received several message packs instead of one")
	ErrorUpdadeMessageStatus = errors.New("failed to update deffered message status")
	ErrorRemoveMessages      = errors.New("failed to execute delete message statement")
)
