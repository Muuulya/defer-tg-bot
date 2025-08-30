package manager

import (
	"fmt"
	"log"

	tgbotapi "github.com/OvyFlash/telegram-bot-api"
)

type MessageManager struct {
	tgAPI *tgbotapi.BotAPI
}

func NewMessageManager(tgAPI *tgbotapi.BotAPI) (sender *MessageManager, rer error) {
	manager := MessageManager{tgAPI: tgAPI}
	return &manager, nil
}

func (manager *MessageManager) SendMessage(targetChatID int64, text string) (message tgbotapi.Message, error error) {
	newMessage := tgbotapi.NewMessage(targetChatID, text)
	message, error = manager.tgAPI.Send(newMessage)
	if error != nil {
		log.Println(error)
		return message, fmt.Errorf("message to chat id %d don't send", targetChatID)
	}

	return message, nil
}

func (manager *MessageManager) SendMessageWithButtons(targetChatID int64, text string, keyboard tgbotapi.ReplyKeyboardMarkup) (message tgbotapi.Message, error error) {
	newMessage := tgbotapi.NewMessage(targetChatID, text)
	newMessage.ReplyMarkup = keyboard

	message, error = manager.tgAPI.Send(newMessage)
	if error != nil {
		log.Println(error)
		return message, fmt.Errorf("message to chat id %d don't send", targetChatID)
	}

	return message, nil
}

func (manager *MessageManager) RemoveMessage(targetChatID int64, messagesID int) {
	deleteConfig := tgbotapi.NewDeleteMessage(targetChatID, messagesID)
	if _, err := manager.tgAPI.Request(deleteConfig); err != nil {
		log.Printf("delete %d: %v", messagesID, err)
	}
}

func (manager *MessageManager) RemoveMessages(targetChatID int64, messagesID []int) {
	const chunkSize = 100

	for start := 0; start < len(messagesID); start += chunkSize {
		end := start + chunkSize
		if end > len(messagesID) {
			end = len(messagesID)
		}

		deleteConfig := tgbotapi.NewDeleteMessages(targetChatID, messagesID[start:end])
		if _, err := manager.tgAPI.Request(deleteConfig); err != nil {
			log.Printf("delete %d-%d: %v", start, end-1, err)
		}

	}
}

func (manager *MessageManager) SendMessageWithInlineButtons(
	targetChatID int64,
	text string,
	keyboard tgbotapi.InlineKeyboardMarkup,
) (message tgbotapi.Message, error error) {
	newMessage := tgbotapi.NewMessage(targetChatID, text)

	newMessage.ReplyMarkup = keyboard

	message, error = manager.tgAPI.Send(newMessage)
	if error != nil {
		log.Println(error)
		return message, fmt.Errorf("message to chat id %d don't send", targetChatID)
	}

	return message, nil
}

func (manager *MessageManager) EditMessage(
	targetChatID int64,
	targetMessageID int,
	text string,
) (message tgbotapi.Message, error error) {
	newMessage := tgbotapi.NewEditMessageText(targetChatID, targetMessageID, text)

	return manager.tgAPI.Send(newMessage)
}

func (manager *MessageManager) EditMessageWithInlineButtons(
	targetChatID int64,
	targetMessageID int,
	text string,
	keyboard tgbotapi.InlineKeyboardMarkup,
) (message tgbotapi.Message, error error) {
	newMessage := tgbotapi.NewEditMessageTextAndMarkup(targetChatID, targetMessageID, text, keyboard)

	return manager.tgAPI.Send(newMessage)
}

func (manager *MessageManager) SendCallbackMessage(callbackQuery *tgbotapi.CallbackQuery, text string) {
	manager.tgAPI.Request(tgbotapi.NewCallback(callbackQuery.ID, text))
}
