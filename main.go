package main

import (
	"context"
	"flag"
	"log"
	tgClient "tgbot2/clients/telegram"
	event_consumer "tgbot2/consumer/event-consumer"
	"tgbot2/events/telegram"
	"tgbot2/files_storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	sqliteStoragePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	s, err := sqlite.New(sqliteStoragePath)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatal("can't init storage: ", err)
	}

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
