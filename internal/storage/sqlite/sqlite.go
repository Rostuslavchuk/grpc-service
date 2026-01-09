package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"sso/internal/domain/models"
	"sso/internal/storage"

	sqlite3 "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s %w", op, err)
	}

	return &Storage{
		db: db,
	}, nil
}

func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users (email, pass_hash) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s %w", op, err)
	}

	sqlResult, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var errSqlite sqlite3.Error
		if errors.As(err, &errSqlite) {
			if errSqlite.ExtendedCode == sqlite3.ErrConstraintUnique {
				return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
			}
		}
		return 0, fmt.Errorf("%s %w", op, err)
	}

	lastInsertedID, err := sqlResult.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return lastInsertedID, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s %w", op, err)
	}

	sqlResult := stmt.QueryRowContext(ctx, email)

	var user models.User
	err = sqlResult.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	stmt, err := s.db.Prepare("SELECT is_admin FROM USERS WHERE id = ?")
	if err != nil {
		return false, fmt.Errorf("%s %w", op, err)
	}

	sqlResult := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = sqlResult.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT * FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s %w", op, err)
	}

	sqlResult := stmt.QueryRowContext(ctx, appID)

	var app models.App
	err = sqlResult.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s %w", op, err)
	}

	return app, nil
}
