package sqlite

import (
	"context"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"tgbot2/files_storage"
	"tgbot2/lib/e"
)

type Storage struct {
	db *sql.DB
}

// New creates new SQLite files_storage.
func New(path string) (*Storage, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, e.WrapIfErr("can't open database", err)
	}

	if err := db.Ping(); err != nil {
		return nil, e.WrapIfErr("can't ping database", err)
	}

	return &Storage{db: db}, nil
}

// Save saves page to files_storage.
func (s *Storage) Save(ctx context.Context, p *files_storage.Page) error {
	q := `INSERT INTO pages (url, user_name) VALUES (?, ?)`

	_, err := s.db.ExecContext(ctx, q, p.URL, p.UserName)
	if err != nil {
		return e.WrapIfErr("can't insert page", err)
	}

	return nil
}

// PickRandom picks random url from files_storage.
func (s *Storage) PickRandom(cxt context.Context, userName string) (*files_storage.Page, error) {
	q := `SELECT url FROM pages WHERE user_name = ? ORDER BY RANDOM() LIMIT 1`

	var url string

	err := s.db.QueryRowContext(cxt, q, userName).Scan(&url)
	if err == sql.ErrNoRows {
		return nil, files_storage.ErrNoSavedPages
	}
	if err != nil {
		return nil, e.WrapIfErr("can't select random page", err)
	}

	return &files_storage.Page{
		URL:      url,
		UserName: userName,
	}, nil
}

// Remove removes url from files_storage.
func (s *Storage) Remove(cxt context.Context, page *files_storage.Page) error {
	q := `DELETE FROM pages WHERE url = ? AND user_name = ?`

	_, err := s.db.ExecContext(cxt, q, page.URL, page.UserName)
	if err != nil {
		return e.WrapIfErr("can't delete page", err)
	}

	return nil
}

// isExists checks if page exists in files_storage.
func (s *Storage) IsExists(cxt context.Context, page *files_storage.Page) (bool, error) {
	q := `SELECT COUNT(*) FROM pages WHERE url = ? AND user_name = ?`

	var count int

	err := s.db.QueryRowContext(cxt, q, page.URL, page.UserName).Scan(&count)
	if err != nil {
		return false, e.WrapIfErr("can't check if page exists", err)
	}

	return count > 0, nil
}

func (s *Storage) Init(cxt context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS pages (url TEXT, user_name TEXT)`
	_, err := s.db.ExecContext(cxt, q)
	if err != nil {
		return e.WrapIfErr("can't create table", err)
	}

	return nil
}
