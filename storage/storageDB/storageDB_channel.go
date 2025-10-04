package storageDB

import (
	"database/sql"
	"log"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/storage"
)

func (s *StorageDB) AddChannel(userID int64, channel *data.Channel) error {
	stmt, err := s.db.Prepare(addChannelSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorAddChannel
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, channel.ID(), channel.Name())
	if err != nil {
		log.Println(err)
		return storage.ErrorAddChannel
	}

	return nil
}

func (s *StorageDB) UpdateChannelName(userID int64, channel *data.Channel) error {
	stmt, err := s.db.Prepare(updateChannelNameSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateChannel
	}
	defer stmt.Close()

	_, err = stmt.Exec(channel.Name(), userID, channel.ID())
	if err != nil {
		log.Println(err)
		return storage.ErrorUpdateChannel
	}

	return nil
}

func (s *StorageDB) readChannel(scanner scanner) (channel *data.Channel, err error) {
	var id int64
	var name string

	if err = scanner.Scan(
		&id,
		&name,
	); err == nil {
		channel = data.NewChannel(id, name)
	}

	return
}

func (s *StorageDB) GetChannel(userID int64, channelID int64) (channel *data.Channel, err error) {
	row := s.db.QueryRow(getChannelSQL, userID, channelID)
	channel, err = s.readChannel(row)
	if err == sql.ErrNoRows {
		return nil, storage.ErrorChannelNotFound
	}
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetChannel
	}

	return
}

func (s *StorageDB) GetAllUserChannels(userID int64) (channels []data.Channel, error error) {
	rows, err := s.db.Query(getUserChannelsSQL, userID)
	if err != nil {
		log.Println(err)
		return nil, storage.ErrorGetChannels
	}
	defer rows.Close()

	for rows.Next() {
		if channel, err := s.readChannel(rows); err != nil {
			log.Println(err)
		} else {
			channels = append(channels, *channel)
		}
	}

	return channels, nil
}

func (s *StorageDB) RemoveChannel(userID int64, channelID int64) error {
	stmt, err := s.db.Prepare(removeChannelSQL)
	if err != nil {
		log.Println(err)
		return storage.ErrorRemoveChannel
	}
	defer stmt.Close()

	_, err = stmt.Exec(userID, channelID)
	if err != nil {
		log.Println(err)
		return storage.ErrorRemoveChannel
	}

	return nil
}
