package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/muuulya/defer-tg-bot/bot"
	"github.com/muuulya/defer-tg-bot/bot/data"

	"github.com/joho/godotenv"
)

func main() {
	env := getENV()
	bot := bot.NewBot(env)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	bot.Start(ctx)
	<-ctx.Done()
	bot.Stop()
}

func getENV() *data.ENV {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	botToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	if botToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required in .env")
	}

	debugMod := false
	debugModString := os.Getenv("TELEGRAM_BOT_DEBUG_MOD")
	debugModString = strings.ToLower(debugModString)
	if debugModString == "true" {
		debugMod = true
	}

	return &data.ENV{BotToken: botToken, DebugMod: debugMod}
}
