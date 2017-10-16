package main

import (
	"alertbot/telegram"
	"alertbot/interfaces"
)

func gotWebhook(bot interfaces.IBot, messages <-chan string) {
	for m := range messages {
		bot.Broadcast(m);
	}
}

func startBot(config Config, messages <-chan string) {
	bot := telegram.TelegramBot{Token: config.Token}

	for channel := range config.Channel {
		bot.Join(config.Channel[channel].Id)
	}

	go gotWebhook(&bot, messages)

	bot.Connect()
	bot.Run()
}
