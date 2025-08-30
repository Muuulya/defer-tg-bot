package data

import "time"

type DefferedMessage struct {
	UserID         int64
	TargetChanalID int64
	MessageID      int
	Datatime       time.Time
}

func NewDummyDefferedMessage() *DefferedMessage {
	return &DefferedMessage{}
}
