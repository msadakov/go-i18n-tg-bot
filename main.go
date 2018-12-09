package main

import (
	"flag"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/msadakov/go-i18n-tg-bot/log"
)

var token string

func init() {
	flag.StringVar(&token, "token", "", "tg bot token")
}

func main() {
	flag.Parse()

	if token == "" {
		log.Fatal("Token is empty\n")
	}

	bot, err := tg.NewBotAPI(token)
	if err != nil {
		log.Fatal(err.Error())
	}

	bot.Debug = true

	log.Info("Authorized on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("Can't get the update channel: %v", err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Info("[%s]: %s", update.Message.From.UserName, update.Message.Text)

		msg := tg.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(msg)
		if err != nil {
			log.Error("Can't send a message: %v", err)
		}
	}
}
