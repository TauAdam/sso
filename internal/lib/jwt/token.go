package jwt

import (
	"github.com/TauAdam/sso/internal/entities/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func NewToken(user models.User, app models.App, ttl time.Duration) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["app_id"] = app.ID
	claims["exp"] = time.Now().Add(ttl).Unix()

	tokenStr, err := token.SignedString([]byte(app.Secret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
