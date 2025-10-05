package storageDB

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/muuulya/defer-tg-bot/bot/data"
	"github.com/muuulya/defer-tg-bot/storage"
)

func newTestStorage(t *testing.T) *StorageDB {
	t.Helper()

	tmpDir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get working dir: %v", err)
	}
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir to temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			panic(err)
		}
	})

	storageDB, err := NewStorageDB()
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	t.Cleanup(func() {
		storageDB.Close()
	})

	return storageDB
}

func countScheduledMessages(t *testing.T, storage *StorageDB) int {
	t.Helper()

	row := storage.db.QueryRow("SELECT COUNT(*) FROM scheduled_posts")
	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatalf("count scheduled posts: %v", err)
	}

	return count
}

func addUserWithChannel(t *testing.T, storage *StorageDB, userID, channelID int64) {
	t.Helper()

	user := data.NewUser(userID, "test user")
	if err := storage.AddUser(user); err != nil {
		t.Fatalf("add user: %v", err)
	}

	channel := data.NewChannel(channelID, "test channel")
	if err := storage.AddChannel(userID, channel); err != nil {
		t.Fatalf("add channel: %v", err)
	}
}

func TestStorageDB_GetUserNotFound(t *testing.T) {
	db := newTestStorage(t)

	if _, err := db.GetUser(42); !errors.Is(err, storage.ErrorUserNotFound) {
		t.Fatalf("expected ErrorUserNotFound, got %v", err)
	}
}

func TestStorageDB_AddMessagePack_RollbackOnPartialFailure(t *testing.T) {
	db := newTestStorage(t)

	const (
		userID    int64 = 1
		channelID int64 = 10
	)
	addUserWithChannel(t, db, userID, channelID)

	postedTime := time.Now().UTC().Truncate(time.Second)
	pack := data.NewDefferedMessagePack(postedTime, 2)

	valid := data.NewDefferedMessage(userID, channelID, 100)
	invalid := data.NewDefferedMessage(userID, channelID+1, 101)

	pack.AddMessage(valid)
	pack.AddMessage(invalid)

	if err := db.AddMessagePack(pack); !errors.Is(err, storage.ErrorAddMessage) {
		t.Fatalf("expected ErrorAddMessage, got %v", err)
	}

	if count := countScheduledMessages(t, db); count != 0 {
		t.Fatalf("expected no scheduled messages, got %d", count)
	}
}

func TestStorageDB_GetMessagePackForUserChannelAfter_MultiplePacksError(t *testing.T) {
	db := newTestStorage(t)

	const (
		userID    int64 = 1
		channelID int64 = 10
	)
	addUserWithChannel(t, db, userID, channelID)

	baseTime := time.Now().UTC().Truncate(time.Second)

	firstPack := data.NewDefferedMessagePack(baseTime, 1)
	firstPack.AddMessage(data.NewDefferedMessage(userID, channelID, 100))
	if err := db.AddMessagePack(firstPack); err != nil {
		t.Fatalf("add first pack: %v", err)
	}

	secondPack := data.NewDefferedMessagePack(baseTime.Add(time.Hour), 1)
	secondPack.AddMessage(data.NewDefferedMessage(userID, channelID, 101))
	if err := db.AddMessagePack(secondPack); err != nil {
		t.Fatalf("add second pack: %v", err)
	}

	_, err := db.GetMessagePackForUserChannelAfter(userID, channelID, baseTime.Add(-time.Minute))
	if !errors.Is(err, storage.ErrorSeveralMessagePacks) {
		t.Fatalf("expected ErrorSeveralMessagePacks, got %v", err)
	}
}

func TestStorageDB_GetMissedMessagesPacksBefor_NoRows(t *testing.T) {
	db := newTestStorage(t)

	if _, err := db.GetMissedMessagesPacksBefor(time.Now()); !errors.Is(err, storage.ErrorMessageNotFound) {
		t.Fatalf("expected ErrorMessageNotFound, got %v", err)
	}
}

func TestStorageDB_UpdateMessagePackStatus(t *testing.T) {
	db := newTestStorage(t)

	const (
		userID    int64 = 1
		channelID int64 = 10
	)
	addUserWithChannel(t, db, userID, channelID)

	scheduledTime := time.Now().UTC().Truncate(time.Second)
	pack := data.NewDefferedMessagePack(scheduledTime, 2)
	pack.AddMessage(data.NewDefferedMessage(userID, channelID, 100))
	pack.AddMessage(data.NewDefferedMessage(userID, channelID, 101))

	if err := db.AddMessagePack(pack); err != nil {
		t.Fatalf("add pack: %v", err)
	}

	storedPack, err := db.GetMessagePackForUserChannelAfter(userID, channelID, scheduledTime.Add(-time.Minute))
	if err != nil {
		t.Fatalf("get pack: %v", err)
	}

	for i := range storedPack.Messages() {
		storedPack.SetMessageStatus(i, true)
	}

	if err := db.UpdateMessagePackStatus(storedPack); err != nil {
		t.Fatalf("update pack status: %v", err)
	}

	rows, err := db.db.Query("SELECT is_posted FROM scheduled_posts ORDER BY id")
	if err != nil {
		t.Fatalf("query scheduled posts: %v", err)
	}
	defer rows.Close()

	index := 0
	for rows.Next() {
		var posted bool
		if err := rows.Scan(&posted); err != nil {
			t.Fatalf("scan post status: %v", err)
		}
		if !posted {
			t.Fatalf("expected message %d to be marked as posted", index)
		}
		index++
	}

	if err := rows.Err(); err != nil {
		t.Fatalf("rows error: %v", err)
	}

	if index != 2 {
		t.Fatalf("expected 2 messages, got %d", index)
	}
}
