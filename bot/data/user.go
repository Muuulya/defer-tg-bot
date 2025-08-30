package data

type User struct {
	ID                     int64
	Name                   string
	CurrentStateName       string
	CurrentDialogMessageID int
	CurrentChannelPage     int
	SelectedChannelID      int64
}

func NewDummyUser() *User {
	return &User{}
}

func NewUser(id int64, name string) *User {
	user := NewDummyUser()
	user.ID = id
	user.Name = name
	user.CurrentDialogMessageID = 0
	user.CurrentChannelPage = 0
	user.SelectedChannelID = 0
	return user
}
