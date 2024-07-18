package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(psqlInfo string) (*Storage, error) {
	const fn = "storage.postgresql.New"
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS url (
			alias VARCHAR(255) PRIMARY KEY,
			url TEXT NOT NULL);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	stmt, err = db.Prepare(`
		CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec()

	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) error {
	const fn = "storage.postgresql.SaveURL"

	stmt, err := s.db.Prepare(`
		INSERT INTO url VALUES($1, $2);
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	_, err = stmt.Exec(alias, urlToSave)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const fn = "storage.postgresql.GetURL"

	stmt, err := s.db.Prepare(`
		SELECT (url) FROM url 
		WHERE alias = $1;
	`)

	if err != nil {
		return "", fmt.Errorf("%s: %w", fn, err)
	}

	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL)

	if err != nil {
		return "", fmt.Errorf("%s: execute statement: %w", fn, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const fn = "storage.postgresql.GetURL"

	stmt, err := s.db.Prepare(`
		DELETE FROM url 
		WHERE alias = $1;
	`)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	res, err := stmt.Exec(alias)

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: Alias does not exist", fn)
	}

	return nil
}
