package jwt

import (
	"testing"
	"time"

	"sso/internal/domain/models"
)

var user = &models.User{
	ID:       2,
	Email:    "rostyk@gmail.com",
	PassHash: []byte("sfvwsewfef"),
}

var app = &models.App{
	ID:     3,
	Name:   "lokeded",
	Secret: "fwfwfwefa wrfawerf",
}

func TestJWT(t *testing.T) {
	testData := []struct {
		Name     string
		user     models.User
		app      models.App
		duration time.Duration
	}{
		{
			Name:     "test 1",
			user:     *user,
			app:      *app,
			duration: 2 * time.Minute,
		},
	}
	for _, tt := range testData {
		t.Run(tt.Name, func(t *testing.T) {
			got, err := NewToken(tt.user, tt.app, tt.duration)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(got)
		})
	}
}
