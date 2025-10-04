package data

type User struct {
	id                     int64
	name                   string
	currentStateName       string
	currentDialogMessageID int
	currentChannelPage     int
	selectedChannelID      int64
}

func NewUser(id int64, name string) *User {
	user := User{}
	user.id = id
	user.name = name
	user.currentDialogMessageID = 0
	user.currentChannelPage = 0
	user.selectedChannelID = 0
	return &user
}

func (u *User) ID() int64                   { return u.id }
func (u *User) Name() string                { return u.name }
func (u *User) CurrentStateName() string    { return u.currentStateName }
func (u *User) CurrentDialogMessageID() int { return u.currentDialogMessageID }
func (u *User) CurrentChannelPage() int     { return u.currentChannelPage }
func (u *User) SelectedChannelID() int64    { return u.selectedChannelID }

func (u *User) SetCurrentState(staneName string) { u.currentStateName = staneName }
func (u *User) SetDialogMessageID(id int)        { u.currentDialogMessageID = id }
func (u *User) SetChannelPage(page int)          { u.currentChannelPage = page }
func (u *User) SetSelectedChannel(id int64)      { u.selectedChannelID = id }
