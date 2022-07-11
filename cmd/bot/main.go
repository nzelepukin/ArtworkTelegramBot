package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	ArtAPI "github.com/nzelepukin/ArtworkTelegramBot/internal/ArtAPI"
	TranslateAPI "github.com/nzelepukin/ArtworkTelegramBot/internal/TranslateAPI"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	TOKEN := getenv("BOT_TOKEN", "")
	DETECT_LANG_URL := getenv("DETECT_LANG_URL", "")
	bot, err := tgbotapi.NewBotAPI(TOKEN)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			artworkquery := update.Message.Text
			if !(DETECT_LANG_URL == "") {
				artworkquery = TranslateAPI.TranslateAPI(update.Message.Text, DETECT_LANG_URL)
			}
			artwork := ArtAPI.GetArtAPI(artworkquery)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, artwork)
			msg.ReplyToMessageID = update.Message.MessageID

			bot.Send(msg)
		}
	}
}
