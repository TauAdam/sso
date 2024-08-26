package auth

import (
	"context"
	"errors"
	"fmt"
	"github.com/TauAdam/sso/internal/entities/models"
	"github.com/TauAdam/sso/internal/lib/jwt"
	"github.com/TauAdam/sso/internal/lib/logger/sl"
	"github.com/TauAdam/sso/internal/services/storage"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Auth struct {
	log         *slog.Logger
	tokenTTL    time.Duration
	userStore   UserStore
	userRecord  UserRecord
	appProvider AppProvider
}

var (
	ErrWrongCredentials = errors.New("wrong credentials")
	ErrInvalidAppID     = errors.New("invalid app ID")
	ErrAlreadyExists    = errors.New("user already exists")
)

// UserStore is the interface for creating new user record.
type UserStore interface {
	StoreUser(
		ctx context.Context,
		email string,
		hashedPass []byte,
	) (uid int64, err error)
}

type UserRecord interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(
	log *slog.Logger,
	tokenTTL time.Duration,
	userStore UserStore,
	userRecord UserRecord,
	appProvider AppProvider,
) *Auth {
	return &Auth{
		log:         log,
		tokenTTL:    tokenTTL,
		userStore:   userStore,
		userRecord:  userRecord,
		appProvider: appProvider,
	}
}

// Login checks the email and password and returns a token if the user is authenticated.
func (a *Auth) Login(ctx context.Context, email, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("logging in")

	user, err := a.userRecord.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrWrongCredentials)
		}

		log.Error("failed to find user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(
		user.HashedPassword,
		[]byte(password),
	); err != nil {
		log.Error("invalid password", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrWrongCredentials)
	}

	app, err := a.appProvider.App(ctx, int64(appID))
	if err != nil {
		log.Error("failed to get app", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		log.Error("failed to create token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user with email logged in", slog.String("email", email))

	return token, nil
}

// RegisterUser registers a new user and returns the user ID.
func (a *Auth) RegisterUser(ctx context.Context, email, password string) (int64, error) {
	const op = "auth.RegisterUser"

	log := a.log.With(slog.String("op", op), slog.String("email", email))

	log.Info("registering user")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		log.Error("failed to hash password", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userStore.StoreUser(ctx, email, hash)
	if err != nil {
		if errors.Is(err, storage.ErrUserDuplicate) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrAlreadyExists)
		}

		log.Error("failed to store user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered", slog.Int64("id", id))

	return id, nil
}

// IsAdmin checks if the user is an admin.
func (a *Auth) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(slog.String("op", op), slog.Int64("userID", userID))

	isAdmin, err := a.userRecord.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			log.Warn("app not found", sl.Err(err))
			return false, fmt.Errorf("%s: %w", op, ErrInvalidAppID)
		}
		log.Error("failed to check if user is admin", sl.Err(err))
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user is admin", slog.Bool("isAdmin", isAdmin))

	return isAdmin, nil
}
