package jwt

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type Token struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
}

type TokenClaims struct {
	Guid string `json:"guid"`
	jwt.StandardClaims
}

var signingKey = []byte("changegoal")

func (t *Token) NewAccessToken(guid string) (string, error) {

	expiresAt := time.Now().Add(time.Minute * 5).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, &TokenClaims{
		guid,
		jwt.StandardClaims{ExpiresAt: expiresAt},
	})

	return token.SignedString(signingKey)
}

func (t *Token) CreateRefreshToken() (string, string, error) {
	var tokenCrypt []byte
	var err error

	b := make([]byte, 10)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err = r.Read(b); err != nil {
		return "", "", err
	}
	if tokenCrypt, err = bcrypt.GenerateFromPassword(b, 14); err == nil {
		tokenStr := base64.StdEncoding.EncodeToString(b)
		return string(tokenCrypt), tokenStr, nil
	}

	return "", "", err
}

func (t *Token) ParseAccessToken(accessToken string) (*TokenClaims, error) {

	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return nil, err
	}

	if token.Valid {
		claims, ok := token.Claims.(*TokenClaims)
		if !ok {
			return nil, errors.New("token claims are not are not if type *TokenClaims")
		}

		return claims, nil
	}

	return nil, fmt.Errorf("couldn't parse access token: %w", err)
}
