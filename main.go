package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/crocone/tg-bot"

	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

var startMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Say hello", "hi"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Say bye", "by"),
	),
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_API_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.CallbackQuery != nil {
			callbacks(update)
		} else if update.Message != nil && update.Message.IsCommand() {
			commands(update)
		} else {
			// simply message
		}
	}
}

func callbacks(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.From.ID

	firstName := update.CallbackQuery.From.FirstName
	lastName := update.CallbackQuery.From.LastName

	var text string
	switch data {
	case "hi":
		text = fmt.Sprintf("hi: %v %v", firstName, lastName)
	case "by":
		text = fmt.Sprintf("bye: %v %v", firstName, lastName)
	default:
		text = "Unknown command"
	}
	msg := tgbotapi.NewMessage(chatId, text)
	sendMessage(msg)
}

func commands(update tgbotapi.Update) {
	command := update.Message.Command()
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose options")
		msg.ReplyMarkup = startMenu
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	}
}

func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Panicf("Send message error: %v", err)
	}
}
