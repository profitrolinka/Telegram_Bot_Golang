package event_consumer

import (
	"log"
	"tgbot2/events"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) Consumer {
	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())

			continue
		}

		/*
			есть 2 метода получения запросов:
			1. webhook
			2. long poll
		*/

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := c.handleEvents(gotEvents); err != nil {
			log.Printf("[ERR] consumer handle events: %s", err.Error())

			continue
		}
	}
}

/*
Проблемы с функцией handleEvents:

1. Потеря событий (

возможно решить:
ретраями(но не 100%),
возвращение в хранилище(не самый надёжный, но простой),
фоллбэк( лучший вариант)
)

2. обработка всей пачки (

возможно решить:
остановка после первой ошибки
остановка после некоторого количества ошибок(чуть лучше)
)

3. параллельная обработка

TODO: Изучить sync.WaitGroup{}
*/

func (c Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
