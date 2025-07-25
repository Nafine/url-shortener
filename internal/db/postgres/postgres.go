package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"url-shortener/internal/db"
)

type Storage struct {
	db *sql.DB
}

func execStmt(db *sql.DB, stmtStr string, args ...any) (sql.Result, error) {
	const op = "db.postgres.execStmt"
	stmt, err := db.Prepare(stmtStr)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var res sql.Result
	if args == nil {
		res, err = stmt.Exec()
	} else {
		res, err = stmt.Exec(args...)
	}

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}

func New(storagePath string) (*Storage, error) {
	const op = "db.postgres.New"

	conn, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = execStmt(conn, `
		CREATE TABLE IF NOT EXISTS url (
		id SERIAL PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL UNIQUE);`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = execStmt(conn, "CREATE INDEX IF NOT EXISTS urls_alias_idx ON url USING hash(alias);")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: conn}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "db.postgres.SaveURL"

	stmt, err := s.db.Prepare(`INSERT INTO url (alias, url) VALUES ($1, $2)`)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(alias, urlToSave)
	if err != nil {
		var pqError *pq.Error
		if errors.As(err, &pqError) && pqError.Code == "23505" { //unique_violation
			return fmt.Errorf("%s: %w", op, db.ErrURLAlreadyExists)
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "db.postgres.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = $1")
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string
	err = stmt.QueryRow(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("%s: %w", op, db.ErrURLNotFound)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return url, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "db.postgres.DeleteURL"

	res, err := execStmt(s.db, "DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if n == 0 {
		return fmt.Errorf("%s: %w", op, db.ErrURLNotFound)
	}

	return nil
}
