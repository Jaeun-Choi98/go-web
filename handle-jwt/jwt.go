package handlejwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWTToken(userID, role string) (string, error) {
	newClaim := CustomClaims{
		UserId: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "myserver.com",
			Audience:  []string{"myapp.com"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 1)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			//NotBefore: jwt.NewNumericDate(time.Now()),
			//ID: "JWT ID",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaim)
	return token.SignedString([]byte("secret_key"))
}

func ValidateJWT(tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invaild token method")
		}
		return []byte("secret_key"), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, jwt.ErrInvalidKeyType
	}

	return claims, nil
}
