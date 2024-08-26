package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/TauAdam/sso/internal/entities/models"
	"github.com/TauAdam/sso/internal/services/storage"
	"github.com/mattn/go-sqlite3"
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

	stmt, err := s.db.Prepare("INSERT INTO users(email, hashed_password) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	result, err := stmt.ExecContext(ctx, email, hashedPass)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserDuplicate)
		}

		return 0, fmt.Errorf("%s: exec statement: %w", op, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, hashed_password FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.HashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, userID)

	var res bool
	err = row.Scan(&res)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return false, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return res, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: prepare statement: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: scan row: %w", op, err)
	}

	return app, nil
}
