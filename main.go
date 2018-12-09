package main

import (
	"bytes"
	"flag"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/text/currency"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/msadakov/go-i18n-tg-bot/log"
)

var token string

func init() {
	flag.StringVar(&token, "token", "", "tg bot token")
}

func init() {
	message.SetString(language.English,
		"Hello, %s!\nMy name is %s. I'm a telegram bot.",
		"Hello, %s!\nMy name is %s. I'm a telegram bot.")
	message.SetString(language.Russian,
		"Hello, %s!\nMy name is %s. I'm a telegram bot.",
		"Привет, %s!\nМеня зовут %s. Я телегерам-бот.")

	message.Set(language.English,
		"There are %d day(s) left until the new year.",
		plural.Selectf(1, "%d",
			plural.One, "There are one day left until the new year!",
			plural.Other, "There are %d days left until the new year."))
	message.Set(language.Russian,
		"There are %d day(s) left until the new year.",
		plural.Selectf(1, "%d",
			"=0", "С Новым годом!",
			"=1", "До нового года остался один день!",
			"<5", "До нового года осталось %d дня.",
			"other", "До нового года осталось %d дней."))

	message.SetString(language.English,
		"Exchange rates",
		"Exchange rates")
	message.SetString(language.Russian,
		"Exchange rates",
		"Курсы валют")
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

func cmdStart(p *message.Printer, usrName, botName string) string {
	return p.Sprintf("Hello, %s!\nMy name is %s. I'm a telegram bot.", usrName, botName)
}

func cmdNewYear(p *message.Printer) string {
	t := time.Now()
	left := 365 - t.YearDay()

	return p.Sprintf("There are %d day(s) left until the new year.", left)
}

type exchange struct {
	from, to currency.Unit
}

var rates = map[exchange]float64{
	{currency.USD, currency.RUB}: 66.42,
	{currency.USD, currency.EUR}: 0.88,
	{currency.RUB, currency.USD}: 0.015,
	{currency.RUB, currency.EUR}: 0.013,
	{currency.EUR, currency.RUB}: 75.60,
	{currency.EUR, currency.USD}: 1.14,
}

func cmdRate(p *message.Printer, usrLng language.Tag) string {
	unit, _ := currency.FromTag(usrLng)

	var buf bytes.Buffer

	p.Fprintf(&buf, "Exchange rates")
	buf.WriteString(":\n")

	for rate, amount := range rates {
		if rate.to == unit {
			p.Fprintf(&buf, "%s: %.2f\n", currency.Symbol(rate.from), amount)
		}
	}

	return buf.String()
}
