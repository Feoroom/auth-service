package tokens

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Access имеет тип JWT и алгорит SHA512
type JWT struct {
	secretKey string
}

func NewJWT(secretKey string) *JWT {
	return &JWT{secretKey: secretKey}
}

func (t *JWT) CreateToken(userID, userIP, email string, duration time.Duration) (string, *UserClaims, error) {
	claims, err := NewUserClaims(userID, userIP, email, duration)
	if err != nil {
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		return "", nil, err
	}

	return tokenStr, claims, nil
}

func (t *JWT) VerifyToken(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("недопустимый метод подписи токена")
		}

		return []byte(t.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("токен недействителен")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, err
	}

	return claims, nil
}

func (t *JWT) CreateRefreshToken(userID, userIP, email string) (string, string, *UserClaims, error) {

	refreshToken, refreshClaims, err := t.CreateToken(userID, userIP, email, 24*time.Hour)
	if err != nil {
		return "", "", nil, err
	}
	rftBase64 := base64.StdEncoding.EncodeToString([]byte(refreshToken))
	rftBase64Hash := sha256.Sum256([]byte(rftBase64))

	refreshTokenHash, err := bcrypt.GenerateFromPassword([]byte(rftBase64Hash[:]), bcrypt.DefaultCost)
	if err != nil {
		return "", "", nil, err
	}

	return string(refreshTokenHash), refreshToken, refreshClaims, err

}
