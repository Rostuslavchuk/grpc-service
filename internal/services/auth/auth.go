package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"sso/internal/jwt"
	"sso/internal/storage"

	"golang.org/x/crypto/bcrypt"

	"sso/internal/domain/models"
	"sso/internal/lib/sl"
)

func HashPassword(log *slog.Logger, password, op string) ([]byte, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Error while hash password", sl.Err(err))
		return nil, fmt.Errorf("%s %w", op, err)
	}
	return hashed, nil
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHash []byte) (int64, error)
}
type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}
type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}
type UserOperation interface {
	UserSaver
	UserProvider
	AppProvider
}
type Auth struct {
	log      *slog.Logger
	storage  UserOperation
	tokenTTL time.Duration
}

func New(log *slog.Logger, storageOprations UserOperation, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storageOprations,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email, password string, appID int64) (string, error) {
	const op = "New.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	user, err := a.storage.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user is not exists", sl.Err(err))
			return "", fmt.Errorf("%s %w", op, storage.ErrInvalidCredentials)
		}

		log.Error("faild to get user", sl.Err(err))
		return "", fmt.Errorf("%s %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if errors.Is(err, storage.ErrPasswordIncorect) {
		log.Error("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s %w", op, storage.ErrInvalidCredentials)
	}

	app, err := a.storage.App(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("app is not exists", sl.Err(err))
			return "", fmt.Errorf("%s %w", op, storage.ErrAppNotFound)
		}
		return "", fmt.Errorf("%s %w", op, err)
	}

	log.Info("user succefully logged in")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("faild to generate token", sl.Err(err))
		return "", fmt.Errorf("%s %w", op, err)
	}

	log.Info("token succefully generated")

	return token, nil
}

func (a *Auth) SaveUser(ctx context.Context, email, password string) (userID int64, err error) {
	const op = "New.Register"

	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	hashed, err := HashPassword(log, password, op)
	if err != nil {
		log.Error("error while hashing password", sl.Err(err))
		return 0, fmt.Errorf("%s %w", op, err)
	}

	id, err := a.storage.SaveUser(ctx, email, hashed)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			log.Error("Error user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s %w", op, storage.ErrUserExists)
		}
		log.Error("Error while save user", sl.Err(err))
		return 0, fmt.Errorf("%s %w", "Error while save user", err)
	}

	log.Info("user succefully registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {
	const op = "New.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("userID", userID),
	)

	isAdmin, err = a.storage.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Error("user do not exists", sl.Err(err))
			return false, fmt.Errorf("%s %w", op, storage.ErrAppNotFound)
		}
		log.Error("error while checking user permission", sl.Err(err))
		return false, fmt.Errorf("%s %w", op, err)
	}

	log.Info("succefully checked user ststus")

	return isAdmin, err
}
