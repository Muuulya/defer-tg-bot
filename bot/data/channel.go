package data

type Channel struct {
	id   int64
	name string
}

func NewChannel(id int64, name string) *Channel {
	channel := Channel{}
	channel.id = id
	channel.name = name
	return &channel
}

func (ch *Channel) ID() int64    { return ch.id }
func (ch *Channel) Name() string { return ch.name }
