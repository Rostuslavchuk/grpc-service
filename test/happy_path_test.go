package test

import (
	"testing"
	"time"

	"sso/test/suit"

	ssov1 "github.com/Rostuslavchuk/sso-protos/gen/go/sso"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID     = 0
	appID          = 1
	secret         = "test-secret"
	passDefaultLen = 10
)

func TestHappyPath(t *testing.T) {
	ctx, sut := suit.New(t)

	email := gofakeit.Email()
	pass := GeneratePass()

	response, err := sut.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, response.GetUserId())

	respLog, err := sut.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: pass,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLog.GetToken()
	require.NotEmpty(t, token)

	tokenJWT, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenJWT.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, response.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"])
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1
	assert.InDelta(t, loginTime.Add(sut.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterTwice(t *testing.T) {
	ctx, sut := suit.New(t)

	email := gofakeit.Email()
	pass := GeneratePass()

	response, err := sut.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, response.GetUserId())

	_, err = sut.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "user already exists")
}

func TestRegisterBadCreds(t *testing.T) {
	ctx, sut := suit.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "Field Password is required",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    GeneratePass(),
			expectedErr: "Field Email is required",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "Field Email is required, Field Password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sut.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func TestLoginBadCreds(t *testing.T) {
	ctx, sut := suit.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int64
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    GeneratePass(),
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "email is required",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    GeneratePass(),
			appID:       appID,
			expectedErr: "invalid email or password",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    GeneratePass(),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := sut.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			t.Error(err)
		})
	}
}

func GeneratePass() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
