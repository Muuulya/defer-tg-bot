package bot

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) GetUserState(userID int64) string {
	var state string
	err := b.DB.QueryRow("SELECT state FROM users WHERE telegram_id = ?", userID).Scan(&state)
	if err == sql.ErrNoRows {
		b.DB.Exec("INSERT INTO users (telegram_id, state) VALUES (?, ?)", userID, "idle")
		return "idle"
	}
	return state
}

func (b *Bot) SetUserState(userID int64, state string) {
	b.DB.Exec("UPDATE users SET state = ? WHERE telegram_id = ?", state, userID)
}

func (b *Bot) SaveChannelName(userID int64, name string) {
	b.DB.Exec("INSERT INTO channels (user_id, name) VALUES (?, ?)", userID, name)
	b.SetUserState(userID, "adding_channel_link")
	b.API.Send(tgbotapi.NewMessage(userID, "Теперь отправьте ссылку на канал"))
}

func (b *Bot) SaveChannelLink(userID int64, link string) {
	var id int64
	err := b.DB.QueryRow("SELECT id FROM channels WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID).Scan(&id)
	if err != nil {
		b.API.Send(tgbotapi.NewMessage(userID, "Ошибка при сохранении канала"))
		return
	}

	chat, err := b.API.GetChat(tgbotapi.ChatInfoConfig{
		ChatConfig: tgbotapi.ChatConfig{
			SuperGroupUsername: link,
		},
	})
	if err != nil {
		b.API.Send(tgbotapi.NewMessage(userID, "Не удалось получить чат, проверьте ссылку"))
		return
	}
	member, err := b.API.GetChatMember(tgbotapi.GetChatMemberConfig{
		ChatConfigWithUser: tgbotapi.ChatConfigWithUser{
			ChatID: chat.ID,
			UserID: b.API.Self.ID,
		},
	})
	if err != nil || member.Status != "administrator" {
		b.API.Send(tgbotapi.NewMessage(userID, "Добавьте бота в канал с правами администратора"))
		return
	}

	b.DB.Exec("UPDATE channels SET telegram_id = ? WHERE id = ?", chat.ID, id)
	b.SetUserState(userID, "idle")
	b.API.Send(tgbotapi.NewMessage(userID, "Канал успешно добавлен!"))
}

func (b *Bot) StartScheduling(userID, channelID int64) {
	b.DB.Exec("UPDATE users SET state = ? WHERE telegram_id = ?", "scheduling_time", userID)
	b.DB.Exec("UPDATE users SET state = 'scheduling_time' WHERE telegram_id = ?", userID)
	b.DB.Exec("INSERT OR REPLACE INTO users (telegram_id, state) VALUES (?, ?)", userID, fmt.Sprintf("scheduling:%d", channelID))
	b.API.Send(tgbotapi.NewMessage(userID, "Введите дату и время публикации (формат: 2006-01-02 15:04)"))
}

func (b *Bot) SaveScheduledTime(userID int64, input string) {
	channelID := b.getChannelIDFromState(userID)
	loc, _ := time.LoadLocation("Europe/Moscow")
	t, err := time.ParseInLocation("2006-01-02 15:04", input, loc)
	if err != nil {
		b.API.Send(tgbotapi.NewMessage(userID, "Неверный формат времени"))
		return
	}

	b.DB.Exec("INSERT INTO scheduled_posts (user_id, channel_id, message_ids, scheduled_time) VALUES (?, ?, ?, ?)", userID, channelID, "[]", t.UTC().Format(time.RFC3339))
	b.SetUserState(userID, "collecting_messages")
	b.API.Send(tgbotapi.NewMessage(userID, "Отправьте сообщения для публикации. Когда закончите — нажмите кнопку 'Готово'"))

	button := tgbotapi.NewInlineKeyboardButtonData("Готово", "done_collecting")
	keyboard := tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{button})
	msg := tgbotapi.NewMessage(userID, "Нажмите, когда закончите отправку сообщений")
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}

func (b *Bot) CollectMessage(userID int64, msg *tgbotapi.Message) {
	var postID int
	err := b.DB.QueryRow("SELECT id FROM scheduled_posts WHERE user_id = ? ORDER BY id DESC LIMIT 1", userID).Scan(&postID)
	if err != nil {
		b.API.Send(tgbotapi.NewMessage(userID, "Ошибка: не найден запланированный пост"))
		return
	}

	var msgIDsJSON string
	_ = b.DB.QueryRow("SELECT message_ids FROM scheduled_posts WHERE id = ?", postID).Scan(&msgIDsJSON)
	var msgIDs []int
	json.Unmarshal([]byte(msgIDsJSON), &msgIDs)

	msgIDs = append(msgIDs, msg.MessageID)
	newJSON, _ := json.Marshal(msgIDs)
	b.DB.Exec("UPDATE scheduled_posts SET message_ids = ? WHERE id = ?", string(newJSON), postID)
}

func (b *Bot) FinalizeMessages(userID int64) {
	b.SetUserState(userID, "idle")
	b.API.Send(tgbotapi.NewMessage(userID, "Сообщения сохранены! Ждите публикации в назначенное время."))
}

func (b *Bot) PublishMessages(postID int, channelID, userID int64, msgIDs []int) {
	for _, msgID := range msgIDs {
		copyMsg := tgbotapi.NewCopyMessage(channelID, userID, msgID)
		_, err := b.API.Send(copyMsg)
		if err != nil {
			b.API.Send(tgbotapi.NewMessage(userID, fmt.Sprintf("Ошибка при публикации сообщения %d", msgID)))
		}
	}

	b.DB.Exec("DELETE FROM scheduled_posts WHERE id = ?", postID)
	b.API.Send(tgbotapi.NewMessage(userID, "Публикация завершена"))
}

func (b *Bot) getChannelIDFromState(userID int64) int64 {
	var state string
	_ = b.DB.QueryRow("SELECT state FROM users WHERE telegram_id = ?", userID).Scan(&state)
	parts := strings.Split(state, ":")
	if len(parts) == 2 {
		channelID, _ := strconv.ParseInt(parts[1], 10, 64)
		return channelID
	}
	return 0
}
