package identity

import (
	"errors"
	"flomart/config"
	"flomart/domain/user"
	"flomart/internal/identity/dto"
	"github.com/golang-jwt/jwt/v5"
)

// TODO TokenManager (CreateToken, ParseToken, RefreshToken)

var ErrInvalidToken = errors.New("invalid token")
var (
	ErrInvalidKey                = errors.New("key is invalid")
	ErrInvalidKeyType            = errors.New("key is of invalid type")
	ErrHashUnavailable           = errors.New("the requested hash function is unavailable")
	ErrTokenMalformed            = errors.New("token is malformed")
	ErrTokenUnverifiable         = errors.New("token is unverifiable")
	ErrTokenSignatureInvalid     = errors.New("token signature is invalid")
	ErrTokenRequiredClaimMissing = errors.New("token is missing required claim")
	ErrTokenInvalidAudience      = errors.New("token has invalid audience")
	ErrTokenExpired              = errors.New("token is expired")
	ErrTokenUsedBeforeIssued     = errors.New("token used before issued")
	ErrTokenInvalidIssuer        = errors.New("token has invalid issuer")
	ErrTokenInvalidSubject       = errors.New("token has invalid subject")
	ErrTokenNotValidYet          = errors.New("token is not valid yet")
	ErrTokenInvalidId            = errors.New("token has invalid id")
	ErrTokenInvalidClaims        = errors.New("token has invalid claims")
	ErrInvalidType               = errors.New("invalid type for claim")
)

type UserClaim struct {
	jwt.RegisteredClaims
	UserID user.ID `json:"user_id"`
	Role   string  `json:"role"`
}

// TODO сделать методом UserClaim
func CreateToken(userID user.ID, role, secret string, expiry *jwt.NumericDate) (string, error) {
	claims := UserClaim{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiry,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func CreateTokens(userID user.ID, role string) (dto.TAccessToken, dto.TRefreshToken, error) {
	cfg := config.New()
	access, err := CreateToken(userID, role, cfg.AccessTokenSecret, cfg.AccessTokenExpiryHour)
	if err != nil {
		return "", "", err
	}
	refresh, err := CreateToken(userID, role, cfg.RefreshTokenSecret, cfg.RefreshTokenExpiryHour)
	if err != nil {
		return "", "", err
	}
	return dto.TAccessToken(access), dto.TRefreshToken(refresh), err
}

func ParseToken(tokenStr, secret string) (*UserClaim, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaim{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*UserClaim)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}
	return claims, nil
}
