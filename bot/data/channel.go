package data

type Channel struct {
	ID   int64
	Name string
}

func NewDummyChannel() *Channel {
	return &Channel{}
}

func NewChannel(id int64, name string) *Channel {
	channel := NewDummyChannel()
	channel.ID = id
	channel.Name = name
	return channel
}
