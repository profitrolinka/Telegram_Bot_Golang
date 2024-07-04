package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"tgbot2/files_storage"
	"tgbot2/lib/e"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s from %s", text, username)

	if isAddCmd(text) {
		return p.savePage(context.Background(), chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(context.Background(), chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID, username)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(ctx context.Context, chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &files_storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(ctx, page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err = p.storage.Save(ctx, page); err != nil {
		return err
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(ctx context.Context, chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: send random", err) }()

	page, err := p.storage.PickRandom(ctx, username)
	//if err != nil && !errors.Is(err, files_storage.ErrNoSavedPages) {
	//	err := p.tg.SendMessage(chatID, msgCantGivePage)
	//	return err
	//}
	if errors.Is(err, files_storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(ctx, page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int, username string) error {
	return p.tg.SendMessage(chatID, fmt.Sprintf("Привет, %s! \n\n %s", username, msgHelp))
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
