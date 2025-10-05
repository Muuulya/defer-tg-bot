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

	tempDir := t.TempDir()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(wd); err != nil {
			t.Errorf("restore wd: %v", err)
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

func TestStorageDB_GetUser_NotFound(t *testing.T) {
	storageDB := newTestStorage(t)

	_, err := storageDB.GetUser(42)
	if !errors.Is(err, storage.ErrorUserNotFound) {
		t.Fatalf("expected ErrorUserNotFound, got %v", err)
	}
}

func TestStorageDB_AddUser_Idempotent(t *testing.T) {
	storageDB := newTestStorage(t)

	user := data.NewUser(1, "alice")
	if err := storageDB.AddUser(user); err != nil {
		t.Fatalf("add user: %v", err)
	}

	duplicate := data.NewUser(1, "bob")
	if err := storageDB.AddUser(duplicate); err != nil {
		t.Fatalf("add duplicate user: %v", err)
	}

	stored, err := storageDB.GetUser(1)
	if err != nil {
		t.Fatalf("get user: %v", err)
	}
	if stored.Name() != "alice" {
		t.Fatalf("expected persisted name 'alice', got %q", stored.Name())
	}
}

func mustAddUserAndChannel(t *testing.T, storageDB *StorageDB, userID int64, channelID int64) {
	t.Helper()

	user := data.NewUser(userID, "user")
	if err := storageDB.AddUser(user); err != nil {
		t.Fatalf("add user: %v", err)
	}

	channel := data.NewChannel(channelID, "channel")
	if err := storageDB.AddChannel(userID, channel); err != nil {
		t.Fatalf("add channel: %v", err)
	}
}

func TestStorageDB_AddMessagePack_RollbackOnConstraintError(t *testing.T) {
	storageDB := newTestStorage(t)
	mustAddUserAndChannel(t, storageDB, 1, 2)

	postedTime := time.Unix(1_000, 0)
	pack := data.NewDefferedMessagePack(postedTime, 2)

	message := data.NewDefferedMessage(1, 2, 100)
	pack.AddMessage(message)

	duplicate := data.NewDefferedMessage(1, 2, 100)
	pack.AddMessage(duplicate)

	err := storageDB.AddMessagePack(pack)
	if !errors.Is(err, storage.ErrorAddMessage) {
		t.Fatalf("expected ErrorAddMessage, got %v", err)
	}

	var count int
	if err := storageDB.db.QueryRow("SELECT COUNT(*) FROM scheduled_posts").Scan(&count); err != nil {
		t.Fatalf("count scheduled_posts: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no messages inserted after rollback, got %d", count)
	}
}

func TestStorageDB_GetMessagePackForUserChannelAfter_MultiplePacks(t *testing.T) {
	storageDB := newTestStorage(t)
	mustAddUserAndChannel(t, storageDB, 1, 2)

	times := []time.Time{time.Unix(2_000, 0), time.Unix(3_000, 0)}
	for i, ts := range times {
		pack := data.NewDefferedMessagePack(ts, 1)
		msg := data.NewDefferedMessage(1, 2, 200+i)
		pack.AddMessage(msg)
		if err := storageDB.AddMessagePack(pack); err != nil {
			t.Fatalf("add message pack %d: %v", i, err)
		}
	}

	_, err := storageDB.GetMessagePackForUserChannelAfter(1, 2, time.Unix(0, 0))
	if !errors.Is(err, storage.ErrorSeveralMessagePacks) {
		t.Fatalf("expected ErrorSeveralMessagePacks, got %v", err)
	}
}

func TestStorageDB_GetMissedMessagesPacksBefor_NoMessages(t *testing.T) {
	storageDB := newTestStorage(t)

	_, err := storageDB.GetMissedMessagesPacksBefor(time.Now())
	if !errors.Is(err, storage.ErrorMessageNotFound) {
		t.Fatalf("expected ErrorMessageNotFound, got %v", err)
	}
}

func TestStorageDB_GetMessagePackForUserChannelAfter_NotFound(t *testing.T) {
	storageDB := newTestStorage(t)
	mustAddUserAndChannel(t, storageDB, 1, 2)

	_, err := storageDB.GetMessagePackForUserChannelAfter(1, 2, time.Now())
	if !errors.Is(err, storage.ErrorMessageNotFound) {
		t.Fatalf("expected ErrorMessageNotFound, got %v", err)
	}
}
