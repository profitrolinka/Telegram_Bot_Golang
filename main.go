package main

import (
	"flag"
	"log"
	tgClient "tgbot2/clients/telegram"
	event_consumer "tgbot2/consumer/event-consumer"
	"tgbot2/events/telegram"
	"tgbot2/storage/files"
)

const (
	tgBotHost   = "api.telegram.org"
	storagePath = "storage"
	batchSize   = 100
)

func main() {
	s := files.New(storagePath)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, mustToken()),
		s,
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}

}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token of bot you're using",
	)
	flag.Parse()

	if *token == "" {
		log.Fatal("no token? 0_o")
	}
	return *token
}
