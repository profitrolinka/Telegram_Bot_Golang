package storage

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgbot2/lib/e"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("no saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (string, error) {

	hash := sha1.New()

	if _, err := io.WriteString(hash, p.URL); err != nil {
		return "", e.Wrap("can't' calculate hash", err)
	}

	if _, err := io.WriteString(hash, p.UserName); err != nil {
		return "", e.Wrap("can't' calculate hash", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
