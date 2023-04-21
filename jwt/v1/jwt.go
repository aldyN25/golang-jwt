package v1

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Token struct {
	Data interface{}
	*jwt.StandardClaims
}

func GenerateToken(secret string, expiredDuration time.Duration, data interface{}) (*string, error) {
	dataToken := &Token{
		Data: data,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiredDuration).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), dataToken)
	result, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, err
	}
	return &result, err
}
