package data

import (
	"time"
)

type DefferedMessage struct {
	id             int
	userID         int64
	targetChanalID int64
	messageID      int
	postedTime     *time.Time
	isPosted       bool
}

func NewDefferedMessage(
	UserID int64,
	TargetChannelID int64,
	MessageID int,
) *DefferedMessage {
	return &DefferedMessage{
		userID:         UserID,
		targetChanalID: TargetChannelID,
		messageID:      MessageID,
		isPosted:       false,
	}
}

func (dm *DefferedMessage) ID() int { return dm.id }

func (dm *DefferedMessage) UserID() int64 { return dm.userID }

func (dm *DefferedMessage) TargetChannelID() int64 { return dm.targetChanalID }

func (dm *DefferedMessage) MessageID() int { return dm.messageID }

func (dm *DefferedMessage) PostedTime() *time.Time { return dm.postedTime }

func (dm *DefferedMessage) IsPosted() bool { return dm.isPosted }

func (dm *DefferedMessage) SetID(id int) { dm.id = id }

func (dm *DefferedMessage) SetPostedTime(t *time.Time) { dm.postedTime = t }

func (dm *DefferedMessage) SetIsPosted(p bool) { dm.isPosted = p }
