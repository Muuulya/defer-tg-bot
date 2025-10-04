package storageDB

import (
	"database/sql"
	"log"
	"sort"
	"time"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/storage"
)

func (s *StorageDB) AddMessage(defferedMessage *data.DefferedMessage) error {
	stmt, err := s.db.Prepare(addMessageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorAddMessage
	}
	defer stmt.Close()

	if err = s.writeDefferedMessage(stmt, defferedMessage); err != nil {
		log.Println(err)
		return storage.ErrorAddMessage
	}

	return nil
}

func (s *StorageDB) AddMessagePack(pack *data.DefferedMessagePack) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return storage.ErrorAddMessage
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(addMessageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorAddMessage
	}
	defer stmt.Close()

	for _, message := range pack.Messages() {
		if err = s.writeDefferedMessage(stmt, &message); err != nil {
			log.Println(err)
			return storage.ErrorAddMessage
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return storage.ErrorAddMessage
	}

	return nil
}

func (s *StorageDB) writeDefferedMessage(stmt *sql.Stmt, message *data.DefferedMessage) error {
	_, err := stmt.Exec(
		message.UserID(),
		message.TargetChannelID(),
		message.MessageID(),
		message.PostedTime().Unix(),
		message.IsPosted(),
	)

	return err
}

func (s *StorageDB) readDefferedMessage(scanner scanner) (unixTime int64, defferedMessage *data.DefferedMessage, err error) {
	var id int
	var userID int64
	var targetChanalID int64
	var messageID int
	var isPosted bool

	if err = scanner.Scan(
		&id,
		&userID,
		&targetChanalID,
		&messageID,
		&unixTime,
		&isPosted,
	); err == nil {
		postedTime := time.Unix(unixTime, 0)
		defferedMessage = data.NewDefferedMessage(userID, targetChanalID, messageID)
		defferedMessage.SetID(id)
		defferedMessage.SetPostedTime(&postedTime)
		defferedMessage.SetIsPosted(isPosted)
	}

	return unixTime, defferedMessage, err
}

func (s *StorageDB) readDefferedMessages(rows *sql.Rows) (packs []data.DefferedMessagePack, err error) {
	messages := make(map[int64][]data.DefferedMessage)

	for rows.Next() {
		if unixTime, message, err := s.readDefferedMessage(rows); err != nil {
			log.Println(err)
		} else {
			messages[unixTime] = append(messages[unixTime], *message)
		}
	}

	if len(messages) == 0 {
		return nil, storage.ErrorMessageNotFound
	}

	keys := make([]int64, 0, len(messages))
	for k := range messages {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	for _, k := range keys {
		if len(messages[k]) > 0 {
			postedTime := messages[k][0].PostedTime()
			pack := data.NewDefferedMessagePack(*postedTime, len(messages[k]))
			pack.AddMessages(messages[k])
			packs = append(packs, *pack)
		}
	}

	return packs, nil
}

func (s *StorageDB) GetMessagePackForUserChannelAfter(userID int64, channelID int64, after time.Time) (pack *data.DefferedMessagePack, error error) {
	rows, err := s.db.Query(getUserChannelMessagesAfterSQL, userID, channelID, after.Unix())
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetMessages
	}
	defer rows.Close()

	packs, err := s.readDefferedMessages(rows)
	if err == storage.ErrorMessageNotFound {
		return nil, err
	} else if err != nil {
		return nil, storage.ErrorGetMessages
	} else if len(packs) > 1 {
		return nil, storage.ErrorSeveralMessagePacks
	}

	return &packs[0], nil
}

func (s *StorageDB) GetMissedMessagesPacksBefor(befor time.Time) (packs []data.DefferedMessagePack, err error) {
	rows, err := s.db.Query(getAllMissedMessagesBeforSQL, befor.Unix(), false)
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetMessages
	}
	defer rows.Close()

	packs, err = s.readDefferedMessages(rows)
	if err != storage.ErrorMessageNotFound {
		err = storage.ErrorGetMessages
	}
	return packs, err
}

func (s *StorageDB) UpdateMessageStatus(message *data.DefferedMessage) error {
	stmt, err := s.db.Prepare(updateMessageStatusSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdadeMessageStatus
	}
	defer stmt.Close()

	if _, err = stmt.Exec(message.IsPosted(), message.ID()); err != nil {
		log.Println(err)
		return storage.ErrorUpdadeMessageStatus
	}

	return nil
}

func (s *StorageDB) UpdateMessagePackStatus(pack *data.DefferedMessagePack) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdadeMessageStatus
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(updateMessageStatusSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdadeMessageStatus
	}
	defer stmt.Close()

	for _, message := range pack.Messages() {
		if _, err = stmt.Exec(message.IsPosted(), message.ID()); err != nil {
			log.Println(err)
			return storage.ErrorUpdadeMessageStatus
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return storage.ErrorUpdadeMessageStatus
	}

	return nil
}

func (s *StorageDB) RemoveMessage(message *data.DefferedMessage) error {
	stmt, err := s.db.Prepare(removeMessageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorRemoveMessages
	}
	defer stmt.Close()

	if _, err := stmt.Exec(message.ID()); err != nil {
		log.Println(err)
		return storage.ErrorRemoveMessages
	}

	return nil
}

func (s *StorageDB) RemoveMessagePack(pack *data.DefferedMessagePack) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println(err)
		return storage.ErrorRemoveMessages
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	stmt, err := tx.Prepare(removeMessageSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorRemoveMessages
	}
	defer stmt.Close()

	for _, message := range pack.Messages() {
		if _, err = stmt.Exec(message.ID()); err != nil {
			log.Println(err)
			return storage.ErrorRemoveMessages
		}
	}

	if err = tx.Commit(); err != nil {
		log.Println(err)
		return storage.ErrorRemoveMessages
	}

	return nil
}
