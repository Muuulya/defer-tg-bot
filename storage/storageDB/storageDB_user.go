package storageDB

import (
	"database/sql"
	"log"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/storage"
)

func (s *StorageDB) AddUser(user *data.User) error {
	stmt, err := s.db.Prepare(addUserSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorAddUser
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID(), user.Name(), user.CurrentStateName(), user.CurrentDialogMessageID(), user.CurrentChannelPage(), user.SelectedChannelID())
	if err != nil {
		log.Println(err)
		return storage.ErrorAddUser
	}

	return nil
}

func (s *StorageDB) UpdateUserState(user *data.User) error {
	stmt, err := s.db.Prepare(updateUserStateSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentStateName(), user.ID())
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}

	return nil
}

func (s *StorageDB) UpdateUserCurrentDialogMessage(user *data.User) error {
	stmt, err := s.db.Prepare(updateUserCurrentDialogMessageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentDialogMessageID(), user.ID())
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}

	return nil
}

func (s *StorageDB) UpdateUserSelectedChannelID(user *data.User) error {
	stmt, err := s.db.Prepare(updateUserSelectedChannelIDSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.SelectedChannelID(), user.ID())
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}

	return nil
}

func (s *StorageDB) UpdateUserCurrentChannelPage(user *data.User) error {
	stmt, err := s.db.Prepare(updateUserCurrentChannelPageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentChannelPage(), user.ID())
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateUser
	}

	return nil
}

func (s *StorageDB) GetUser(userID int64) (user *data.User, err error) {
	row := s.db.QueryRow(getUserSQL, userID)
	user, err = s.readUser(row)
	if err == sql.ErrNoRows {
		return nil, storage.ErrorUserNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetUser
	}

	return user, nil
}

func (s *StorageDB) GetAllUsers() (users []data.User, err error) {
	rows, err := s.db.Query(getAllUsersSQL)
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetAllUsers
	}
	defer rows.Close()

	users = []data.User{}

	for rows.Next() {
		user, err := s.readUser(rows)
		if err != nil {
			log.Println(err)
		}
		users = append(users, *user)
	}

	return users, nil
}

func (s *StorageDB) readUser(scanner scanner) (user *data.User, err error) {
	var id int64
	var name string
	var currentStateName string
	var currentDialogMessageID int
	var currentChannelPage int
	var selectedChannelID int64

	if err = scanner.Scan(
		&id,
		&name,
		&currentStateName,
		&currentDialogMessageID,
		&currentChannelPage,
		&selectedChannelID,
	); err == nil {
		user = data.NewUser(id, name)
		user.SetCurrentState(currentStateName)
		user.SetDialogMessageID(currentDialogMessageID)
		user.SetChannelPage(currentChannelPage)
		user.SetSelectedChannel(selectedChannelID)
	}

	return
}
