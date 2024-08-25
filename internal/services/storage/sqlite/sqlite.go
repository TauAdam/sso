package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TauAdam/sso/internal/services/storage"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	const op = "sqlite.New"

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("%s: open database: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) StoreUser(ctx context.Context, email string, hashedPass []byte) (int64, error) {
	const op = "sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, email, hashedPass)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.Extended == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserDuplicate)
		}

		return 0, fmt.Errorf("%s: exec statement: %w", op, err)
	}
}
