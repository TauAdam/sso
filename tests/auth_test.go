package tests

import (
	authv1 "github.com/TauAdam/sso/contracts/gen/go/sso"
	"github.com/TauAdam/sso/tests/kit"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	secret = "secret-string"
	appID  = 1

	passwordLength = 8
)

func TestRegisterLogin_HappyPath(t *testing.T) {
	ctx, suite := kit.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passwordLength)

	registerResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	// Check that the registration was successful
	require.NoError(t, err)

	// Check that the user ID is not empty
	assert.NotEmpty(t, registerResp.GetUserId())

	loginResp, err := suite.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	// Check that the login was successful
	require.NoError(t, err)

	loginTimestamp := time.Now()

	tokenStr := loginResp.GetToken()
	// Check that the token is not empty
	require.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	// Check that the token is valid
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	// Check that the claims are correct
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, registerResp.GetUserId(), int64(claims["uid"].(float64)))

	// Check that the token is not expired
	expectedTimestamp := loginTimestamp.Add(suite.Config.TokenTTL).Unix()
	assert.InDelta(t, expectedTimestamp, claims["exp"].(float64), 1)
}

func TestAlreadyRegistered(t *testing.T) {
	ctx, suite := kit.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passwordLength)

	registerResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registerResp.GetUserId())

	secondRegisterResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	// Check that the registration failed
	require.Error(t, err)
	assert.Empty(t, secondRegisterResp.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestLoginAfterRegistration(t *testing.T) {
	ctx, suite := kit.New(t)

	email := gofakeit.Email()
	password := gofakeit.Password(true, true, true, true, false, passwordLength)

	registerResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.NotEmpty(t, registerResp.GetUserId())

	secondRegisterResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	// Check that the registration failed
	require.Error(t, err)
	assert.Empty(t, secondRegisterResp.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

	loginResp, err := suite.AuthClient.Login(ctx, &authv1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID,
	})
	// Check that the login was successful
	require.NoError(t, err)

	loginTimestamp := time.Now()

	tokenStr := loginResp.GetToken()
	require.NotEmpty(t, tokenStr)

	token, err := jwt.Parse(tokenStr, func(tkn *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := token.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, registerResp.GetUserId(), int64(claims["uid"].(float64)))

	expectedTimestamp := loginTimestamp.Add(suite.Config.TokenTTL).Unix()
	assert.InDelta(t, expectedTimestamp, claims["exp"].(float64), 1)
}

func TestFailedRegistration(t *testing.T) {
	ctx, suite := kit.New(t)

	tests := []struct {
		name          string
		email         string
		password      string
		expectedError string
	}{{
		name:          "empty email",
		email:         "",
		password:      gofakeit.Password(true, true, true, true, false, passwordLength),
		expectedError: "email is required",
	},
		{
			name:          "empty password",
			email:         gofakeit.Email(),
			password:      "",
			expectedError: "password is required",
		},
		{
			name:          "empty email and password",
			email:         "",
			password:      "",
			expectedError: "email is required",
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})

			require.Error(t, err)
			assert.Empty(t, registerResp.GetUserId())
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}

func TestFailedLogin(t *testing.T) {
	ctx, suite := kit.New(t)

	tests := []struct {
		name          string
		email         string
		password      string
		appID         int32
		expectedError string
	}{{
		name:          "empty email",
		email:         "",
		password:      gofakeit.Password(true, true, true, true, false, passwordLength),
		appID:         appID,
		expectedError: "email is required",
	},
		{
			name:          "empty password",
			email:         gofakeit.Email(),
			password:      "",
			appID:         appID,
			expectedError: "password is required",
		},
		{
			name:          "empty email and password",
			email:         "",
			password:      "",
			appID:         appID,
			expectedError: "email is required",
		}, {
			name:          "empty app id",
			email:         gofakeit.Email(),
			password:      gofakeit.Password(true, true, true, true, false, passwordLength),
			appID:         0,
			expectedError: "app_id is required",
		}, {
			name:          "wrong credentials",
			email:         gofakeit.Email(),
			password:      gofakeit.Password(true, true, true, true, false, passwordLength),
			appID:         appID,
			expectedError: "wrong email or password",
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registerResp, err := suite.AuthClient.Register(ctx, &authv1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: gofakeit.Password(true, true, true, true, false, passwordLength),
			})
			require.NoError(t, err)
			require.NotEmpty(t, registerResp.GetUserId())

			loginResp, err := suite.AuthClient.Login(ctx, &authv1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			assert.Empty(t, loginResp.GetToken())
			assert.ErrorContains(t, err, tt.expectedError)
		})
	}
}
