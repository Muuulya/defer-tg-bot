package main

import (
	"log"
	"os"

	"defer-tg-bot/bot"
	"defer-tg-bot/db"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required in .env")
	}

	// Init DB
	database, err := db.InitDB("bot.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Init Telegram Bot
	tgBot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatal(err)
	}
	tgBot.Debug = true

	log.Printf("Authorized on account %s", tgBot.Self.UserName)

	// Start Bot
	b := bot.NewBot(tgBot, database)
	b.Start()
}
