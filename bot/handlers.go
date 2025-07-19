package bot

import (
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (b *Bot) HandleMessage(msg *tgbotapi.Message) {
	userID := msg.From.ID

	// Get user state from DB
	state := b.GetUserState(userID)

	switch state {
	case "idle":
		b.ShowMainMenu(msg.Chat.ID)
	case "adding_channel_name":
		b.SaveChannelName(userID, msg.Text)
	case "adding_channel_link":
		b.SaveChannelLink(userID, msg.Text)
	case "scheduling_time":
		b.SaveScheduledTime(userID, msg.Text)
	case "collecting_messages":
		b.CollectMessage(userID, msg)
	default:
		b.ShowMainMenu(msg.Chat.ID)
	}
}

func (b *Bot) HandleCallback(cb *tgbotapi.CallbackQuery) {
	data := cb.Data
	userID := cb.From.ID

	if data == "done_collecting" {
		b.FinalizeMessages(userID)
		return
	}

	if strings.HasPrefix(data, "select_channel:") {
		channelIDStr := strings.TrimPrefix(data, "select_channel:")
		channelID, _ := strconv.Atoi(channelIDStr)
		b.StartScheduling(userID, int64(channelID))
		return
	}
}

func (b *Bot) ShowMainMenu(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, "Что вы хотите сделать?")
	buttons := []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Добавить канал", "add_channel"),
		tgbotapi.NewInlineKeyboardButtonData("Запланировать публикацию", "schedule_post"),
	}
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttons)
	msg.ReplyMarkup = keyboard
	b.API.Send(msg)
}
