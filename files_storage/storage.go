package files_storage

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"tgbot2/lib/e"
)

type Storage interface {
	Save(ctx context.Context, p *Page) error
	PickRandom(ctx context.Context, userName string) (*Page, error)
	Remove(ctx context.Context, p *Page) error
	IsExists(ctx context.Context, p *Page) (bool, error)
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
