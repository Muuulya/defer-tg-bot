package data

import "time"

type DefferedMessagePack struct {
	postedTime time.Time
	messages   []DefferedMessage
}

func NewDefferedMessagePack(PostedTime time.Time, capacity int) *DefferedMessagePack {
	return &DefferedMessagePack{
		postedTime: PostedTime,
		messages:   make([]DefferedMessage, 0, capacity),
	}
}

func (dmp *DefferedMessagePack) PostedTime() time.Time { return dmp.postedTime }

func (dmp *DefferedMessagePack) Messages() []DefferedMessage {
	return append([]DefferedMessage(nil), dmp.messages...)
}

func (dmp *DefferedMessagePack) AddMessage(message *DefferedMessage) {
	message.SetPostedTime(&dmp.postedTime)
	dmp.messages = append(dmp.messages, *message)
}

func (dmp *DefferedMessagePack) AddMessages(messages []DefferedMessage) {
	for i := range messages {
		messages[i].SetPostedTime(&dmp.postedTime)
	}
	dmp.messages = append(dmp.messages, messages...)
}

func (dmp *DefferedMessagePack) SetMessageStatus(index int, status bool) {
	dmp.messages[index].SetIsPosted(status)
}
