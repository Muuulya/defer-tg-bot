package data

type UserChannelPair struct {
	userID    int64
	channelID int64
}

func NewUserChannelPair(userID int64, channelID int64) UserChannelPair {
	return UserChannelPair{userID: userID, channelID: channelID}
}

func (pair *UserChannelPair) UserID() int64    { return pair.userID }
func (pair *UserChannelPair) ChannelID() int64 { return pair.channelID }
