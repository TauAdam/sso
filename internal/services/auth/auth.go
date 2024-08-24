package auth

import (
	"context"
	"fmt"
	"github.com/TauAdam/sso/internal/entities/models"
	"github.com/TauAdam/sso/internal/lib/logger/sl"
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

// UserStore is the interface for creating new user record.
type UserStore interface {
	StoreUser(
		ctx context.Context,
		email string,
		passwordHash []byte,
	) (uid int64, err error)
}

type UserRecord interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (models.App, error)
}

func New(log *slog.Logger, tokenTTL time.Duration, userStore UserStore, userRecord UserRecord, appProvider AppProvider) *Auth {
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
	panic("not implemented")
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
		log.Error("failed to store user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("user registered", slog.Int64("id", id))

	return id, nil
}

// IsAdmin checks if the user is an admin.
func IsAdmin(ctx context.Context, userID int64) (bool, error) {
	panic("not implemented")
}
