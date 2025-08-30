package data

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

type ENV struct {
	BotToken string
	DebugMod bool
}

// type Buttons = []string

type MessagesID []int

// запись в БД
func (messages MessagesID) Value() (driver.Value, error) {
	if messages == nil {
		return "[]", nil
	}
	b, err := json.Marshal(messages)
	return string(b), err // SQLite TEXT
}

// чтение из БД
func (messages *MessagesID) Scan(val interface{}) error {
	switch v := val.(type) {
	case string:
		return json.Unmarshal([]byte(v), messages)
	case []byte:
		return json.Unmarshal(v, messages)
	case nil:
		*messages = MessagesID{}
		return nil
	default:
		return fmt.Errorf("unsupported type %T", v)
	}
}

func (messages *MessagesID) Clear() {
	*messages = MessagesID{}
}
