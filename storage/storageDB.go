package storage

import (
	"database/sql"
	"errors"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/muuulya/defer-tg-bot/bot/data"
)

const (
	fileNameDB = "storage.db"
)

const (
	createStorageErrorMessage                  = "failed to create storage"
	initStorageErrorMessage                    = "failed to init storage"
	addUserErrorMessage                        = "failed to add a user"
	updateUserStateErrorMessage                = "failed to update user state"
	updateUserCurrentDialogMessageErrorMessage = "failed to update user current dialog message id"
	updateUserSelectedChannelIDErrorMessage    = "failed to update user selected channel id"
	updateUserCurrentChannelPageErrorMessage   = "failed to update user current channel page"
	getUserErrorMessage                        = "failed to get user"
	addChannelErrorMessage                     = "failed to add a channel"
	updateChannelNameErrorMessage              = "failed to update a channel name"
	getChannalErrorMessage                     = "failed to get channel"
	getChannalsErrorMessage                    = "failed to get channels"
	removeChannelErrorMessage                  = "failed to delete channel"
	addMessageErrorMessage                     = "failed to add message"
	getMessageErrorMessage                     = "failed to get messages"
	removeMessageErrorMessage                  = "failed to execute delete message statement"
)

const (
	createTablsSQL = `
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY,
			user_name TEXT,
			current_state_name TEXT,
			current_dialog_message_id INTEGER,
			current_channel_page INTAGER,
			selected_channel_id INTEGER
		);

		CREATE TABLE IF NOT EXISTS channels (
			user_id INTEGER,
			channel_id INTEGER,
			channel_name TEXT,
			PRIMARY KEY (user_id, channel_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS scheduled_posts (
			id INTEGER PRIMARY KEY,
			user_id INTEGER,
			channel_id INTEGER,
			message_ids TEXT,
			scheduled_time DATETIME
		);

		CREATE UNIQUE INDEX IF NOT EXISTS idx_unique_post
		ON scheduled_posts(user_id, channel_id, message_ids, scheduled_time);
	`

	addUserSQL = `INSERT OR IGNORE INTO users(user_id, user_name, current_state_name, current_dialog_message_id, current_channel_page, selected_channel_id) 
		VALUES(?, ?, ?, ?, ?, ?)`
	updateUserStateSQL                = `UPDATE users SET current_state_name = ? WHERE user_id = ?`
	updateUserCurrentDialogMessageSQL = `UPDATE users SET current_dialog_message_id = ? WHERE user_id = ?`
	updateUserSelectedChannelIDSQL    = `UPDATE users SET selected_channel_id = ? WHERE user_id = ?`
	updateUserCurrentChannelPageSQL   = `UPDATE users SET current_channel_page = ? WHERE user_id = ?`
	getUserSQL                        = `SELECT * FROM users WHERE user_id = ? LIMIT 1`
	addChannelSQL                     = `INSERT OR IGNORE INTO channels(user_id, channel_id, channel_name)
		VALUES (?, ?, ?);`
	updateChannelNameSQL = `UPDATE channels SET channel_name = ? WHERE user_id = ? AND channel_id = ?`
	getChannelSQL        = `SELECT channel_id, channel_name FROM channels WHERE user_id = ? AND channel_id = ? LIMIT 1`
	getChannelsSQL       = `SELECT channel_id, channel_name FROM channels WHERE user_id = ?`
	removeChannelSQL     = `DELETE FROM channels WHERE user_id = ? AND channel_id = ?`
	addMessageSQL        = `INSERT OR IGNORE INTO scheduled_posts(user_id, channel_id, message_ids, scheduled_time)
		VALUES (?, ?, ?, ?);`
	getAllMessagesSQL  = `SELECT user_id, channel_id, message_ids, scheduled_time FROM scheduled_posts`
	getNextMessagesSQL = `
		SELECT user_id, channel_id, message_ids, scheduled_time
		FROM scheduled_posts
		WHERE scheduled_time = (
			SELECT MIN(scheduled_time)
			FROM scheduled_posts
			WHERE scheduled_time > ?
		);
	`
	removeMessageSQL = `
		DELETE FROM scheduled_posts
		WHERE user_id = ? AND channel_id = ? AND message_ids = ? AND scheduled_time = ?
	`
)

type StorageDB struct {
	db *sql.DB
}

func NewStorageDB() (*StorageDB, error) {
	db, err := sql.Open("sqlite3", fileNameDB)
	if err != nil {
		log.Fatal(err)
		return nil, errors.New(createStorageErrorMessage)
	}

	storage := StorageDB{db: db}
	err = storage.initDB()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &storage, nil
}

func (storage *StorageDB) Close() { storage.db.Close() }

func (storage *StorageDB) AddUser(user *data.User) error {
	stmt, err := storage.db.Prepare(addUserSQL)
	if err != nil {
		log.Println(err)
		return errors.New(addUserErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.ID, user.Name, user.CurrentStateName, user.CurrentDialogMessageID, user.CurrentChannelPage, user.SelectedChannelID)
	if err != nil {
		log.Println(err)
		return errors.New(addUserErrorMessage)
	}

	return nil
}

func (storage *StorageDB) UpdateUserState(user *data.User) error {
	stmt, err := storage.db.Prepare(updateUserStateSQL)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserStateErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentStateName, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserStateErrorMessage)
	}

	return nil
}

func (storage *StorageDB) UpdateUserCurrentDialogMessage(user *data.User) error {
	stmt, err := storage.db.Prepare(updateUserCurrentDialogMessageSQL)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserCurrentDialogMessageErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentDialogMessageID, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserCurrentDialogMessageErrorMessage)
	}

	return nil
}

func (storage *StorageDB) UpdateUserSelectedChannelID(user *data.User) error {
	stmt, err := storage.db.Prepare(updateUserSelectedChannelIDSQL)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserSelectedChannelIDErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.SelectedChannelID, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserSelectedChannelIDErrorMessage)
	}

	return nil
}

func (storage *StorageDB) UpdateUserCurrentChannelPage(user *data.User) error {
	stmt, err := storage.db.Prepare(updateUserCurrentChannelPageSQL)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserCurrentChannelPageErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.CurrentChannelPage, user.ID)
	if err != nil {
		log.Println(err)
		return errors.New(updateUserCurrentChannelPageErrorMessage)
	}

	return nil
}

func (storage *StorageDB) TryGetUser(userID int64) (user *data.User, found bool, error error) {

	user = data.NewDummyUser()

	err := storage.db.QueryRow(getUserSQL, userID).Scan(&user.ID, &user.Name, &user.CurrentStateName, &user.CurrentDialogMessageID, &user.CurrentChannelPage, &user.SelectedChannelID)

	if err == sql.ErrNoRows {
		return &data.User{}, false, nil
	}
	if err != nil {
		log.Println(err)
		return &data.User{}, false, errors.New(getUserErrorMessage)
	}

	return user, true, nil
}

func (storage *StorageDB) AddChannel(userID int64, channel *data.Channel) error {
	stmt, err := storage.db.Prepare(addChannelSQL)

	if err != nil {
		log.Println(err)
		return errors.New(addChannelErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, channel.ID, channel.Name)
	if err != nil {
		log.Println(err)
		return errors.New(addChannelErrorMessage)
	}

	return nil
}

func (storage *StorageDB) UpdateChannelName(userID int64, channel *data.Channel) error {
	stmt, err := storage.db.Prepare(updateChannelNameSQL)

	if err != nil {
		log.Println(err)
		return errors.New(updateChannelNameErrorMessage)
	}

	_, err = stmt.Exec(channel.Name, userID, channel.ID)
	if err != nil {
		log.Println(err)
		return errors.New(updateChannelNameErrorMessage)
	}

	return nil
}

func (storage *StorageDB) TryGetChannel(userID int64, channelID int64) (channel *data.Channel, found bool, error error) {
	channel = data.NewDummyChannel()

	err := storage.db.QueryRow(getChannelSQL, userID, channelID).Scan(&channel.ID, &channel.Name)

	if err == sql.ErrNoRows {
		return &data.Channel{}, false, nil
	}
	if err != nil {
		log.Println(err)
		return &data.Channel{}, false, errors.New(getChannalErrorMessage)
	}

	return channel, true, nil
}

func (storage *StorageDB) GetAllChannels(userID int64) (channels []data.Channel, error error) {
	rows, err := storage.db.Query(getChannelsSQL, userID)
	if err != nil {
		log.Println(err)
		return nil, errors.New(getChannalsErrorMessage)
	}
	defer rows.Close()

	channels = []data.Channel{}

	for rows.Next() {
		channel := data.NewDummyChannel()
		if err := rows.Scan(&channel.ID, &channel.Name); err != nil {
			log.Println(err)
			return nil, errors.New(getChannalErrorMessage)
		}
		channels = append(channels, *channel)
	}

	return channels, nil
}

func (storage *StorageDB) RemoveChannel(userID int64, channelID int64) error {
	stmt, err := storage.db.Prepare(removeChannelSQL)
	if err != nil {
		log.Println(err)
		return errors.New(removeChannelErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, channelID)
	if err != nil {
		log.Println(err)
		return errors.New(removeChannelErrorMessage)
	}

	return nil
}

func (storage *StorageDB) AddDefferedMessage(defferedMessage *data.DefferedMessage) error {
	stmt, err := storage.db.Prepare(addMessageSQL)
	if err != nil {
		log.Println(err)
		return errors.New(addMessageErrorMessage)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		defferedMessage.UserID,
		defferedMessage.TargetChanalID,
		defferedMessage.MessageID,
		defferedMessage.Datatime,
	)
	if err != nil {
		log.Println(err)
		return errors.New(addMessageErrorMessage)
	}

	return nil
}

func (storage *StorageDB) AddDefferedMessages(defferedMessages []data.DefferedMessage) error {
	stmt, err := storage.db.Prepare(addMessageSQL)
	if err != nil {
		log.Println(err)
		return errors.New(addMessageErrorMessage)
	}
	defer stmt.Close()

	for _, message := range defferedMessages {
		_, err = stmt.Exec(
			message.UserID,
			message.TargetChanalID,
			message.MessageID,
			message.Datatime,
		)
		if err != nil {
			log.Println(err)
			return errors.New(addMessageErrorMessage)
		}
	}

	return nil
}

func (storage *StorageDB) GetAllDefferedMessages() (defferedMessages []data.DefferedMessage, error error) {
	rows, err := storage.db.Query(getAllMessagesSQL)
	if err != nil {
		log.Println(err)
		return nil, errors.New(getMessageErrorMessage)
	}
	defer rows.Close()

	defferedMessages = []data.DefferedMessage{}

	for rows.Next() {
		defferedMessage := data.NewDummyDefferedMessage()
		err := rows.Scan(
			&defferedMessage.UserID,
			&defferedMessage.TargetChanalID,
			&defferedMessage.MessageID,
			&defferedMessage.Datatime,
		)
		if err != nil {
			log.Println(err)
			return nil, errors.New(getMessageErrorMessage)
		}
		defferedMessages = append(defferedMessages, *defferedMessage)
	}

	return defferedMessages, nil
}

func (storage *StorageDB) GetNextDefferedMessages(referenceTime *time.Time) (defferedMessages []data.DefferedMessage, error error) {
	rows, err := storage.db.Query(getNextMessagesSQL, referenceTime)
	if err != nil {
		log.Println(err)
		return nil, errors.New(getMessageErrorMessage)
	}
	defer rows.Close()

	defferedMessages = []data.DefferedMessage{}

	for rows.Next() {
		defferedMessage := data.NewDummyDefferedMessage()
		err := rows.Scan(
			&defferedMessage.UserID,
			&defferedMessage.TargetChanalID,
			&defferedMessage.MessageID,
			&defferedMessage.Datatime,
		)
		if err != nil {
			log.Println(err)
			return nil, errors.New(getMessageErrorMessage)
		}
		defferedMessages = append(defferedMessages, *defferedMessage)
	}

	return defferedMessages, nil
}

func (storage *StorageDB) RemoveDefferedMessages(defferedMessages []data.DefferedMessage) error {
	stmt, err := storage.db.Prepare(removeMessageSQL)
	if err != nil {
		log.Println(err)
		return errors.New(removeMessageErrorMessage)
	}
	defer stmt.Close()

	for _, message := range defferedMessages {
		_, err := stmt.Exec(
			message.UserID,
			message.TargetChanalID,
			message.MessageID,
			message.Datatime,
		)
		if err != nil {
			log.Println(err)
			return errors.New(removeMessageErrorMessage)
		}
	}

	return nil
}

func (storage *StorageDB) initDB() error {
	_, err := storage.db.Exec(createTablsSQL)
	if err != nil {
		log.Println(err)
		return errors.New(initStorageErrorMessage)
	}

	return nil
}
