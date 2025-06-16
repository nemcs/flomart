package identity

import (
	"flomart/config"
	"flomart/domain/user"
	"github.com/golang-jwt/jwt/v5"
)

// TODO TokenManager (CreateToken, ParseToken, RefreshToken)

func CreateToken(userID user.ID, role string) (string, error) {
	cfg := config.New()

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     cfg.JWTExpiration,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

func ParseToken(tokenStr, secret string) (jwt.Claims, error) {

	return nil, nil
}
