package telegram

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"

	"strconv"
)

type TelegramBot struct {
	Token string
	bot   *tgbotapi.BotAPI
	channels  []int64
}

func (t *TelegramBot) Connect() {
	bot, err := tgbotapi.NewBotAPI(t.Token)
	if err != nil {
		log.Panic(err)
	}

	t.bot = bot
}

func (t *TelegramBot) Join(id string) {
	id64, err := strconv.ParseInt(id, 10, 64);
	if err != nil {
		log.Panic(err)
	}
	t.channels = append(t.channels, id64)
	log.Printf("Now talking to %s", id)
}

func (t *TelegramBot) Broadcast(message string) {
	for channel := range t.channels {
		t.send(int64(t.channels[channel]), message)
	}
}

func (t *TelegramBot) Run() {
	log.Printf("Authorized on account %s", t.bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := t.bot.GetUpdatesChan(u)
	if err != nil {
		log.Panic(err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		t.send(update.Message.Chat.ID, "Sorry, no commands are implemented")
	}
}

func (t *TelegramBot) send(id int64, message string) {
	msg := tgbotapi.NewMessage(id, message)
	t.bot.Send(msg)
}
