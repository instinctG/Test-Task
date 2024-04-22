package handler

import (
	"context"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"github.com/instinctG/Test-task/internal/db"
	jwT "github.com/instinctG/Test-task/internal/jwt"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
)

type Guid struct {
	Guid string `json:"guid"`
}
type TokenService interface {
	PostRefreshToken(ctx context.Context, guid, refresh string) error
	ReadRefreshToken(ctx context.Context, guid string) (*db.RefreshToken, error)
	UpdateRefreshToken(ctx context.Context, guid, refresh string) error
}

// PostRefreshToken - получает guid с запроса и возвращает новый сгенерированный access и refresh токенов
func (h *Handler) PostRefreshToken(c *gin.Context) {
	//get the guid
	guid := c.Param("id")

	if !IsValidUUID(guid) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "guid is invalid"})
		return
	}

	sendTokenResponse(c, guid, h.Service.PostRefreshToken)
}

// UpdateRefreshToken - получает access и refresh токены с JSON body и передает их по guid в базу данных MongoDB
// для обновления их в базе и выдает пару обновленных access и refresh токенов
func (h *Handler) UpdateRefreshToken(c *gin.Context) {
	guid := c.Param("id")
	if !IsValidUUID(guid) {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "guid is invalid"})
		return
	}
	var token jwT.Token
	var claims *jwT.TokenClaims
	var err error

	// получает access и refresh токен с JSON body
	if err := c.BindJSON(&token); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "couldn't decode access and refresh tokens from JSON body"})
		return
	}

	//проверяет токены на пустую строку
	if token.Access == "" || token.Refresh == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "missing tokens"})
		return
	}

	//парсим access токен и проверяем на валидность access токен
	if claims, err = token.ParseAccessToken(token.Access); claims == nil || err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "error in validating access token"})
		return
	}
	if claims.Guid != guid {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "id doesn't match to the guid in claims"})
		return
	}

	//проверяем на валидность refresh токен
	if err = h.ValidateRefreshToken(c, claims.Guid, token.Refresh); err != nil {
		c.IndentedJSON(http.StatusUnauthorized, gin.H{"error": "couldn't validate refresh token"})
		return
	}

	sendTokenResponse(c, claims.Guid, h.Service.UpdateRefreshToken)
}

func sendTokenResponse(c *gin.Context, guid string, fn func(context.Context, string, string) error) {
	var token jwT.Token

	//создаем access токен
	access, err := token.NewAccessToken(guid)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "couldn't generate a new access token"})
		return
	}

	// создаем refresh токен в виде base64 и в виде bcrypt хеша saveDB сохраняем в БД
	saveDB, refresh, err := token.CreateRefreshToken()
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "couldn't generate refresh token"})
		return
	}

	if err = fn(c, guid, saveDB); err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "couldn't fetch refresh token"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"access": access, "refresh": refresh, "guid": guid})
}

// ValidateRefreshToken - проверка refresh токена на валидность
func (h *Handler) ValidateRefreshToken(ctx context.Context, guid, refresh string) error {
	var err error
	var refreshToken *db.RefreshToken
	var decodeRefresh []byte
	if refreshToken, err = h.Service.ReadRefreshToken(ctx, guid); err == nil {
		if decodeRefresh, err = base64.StdEncoding.DecodeString(refresh); err == nil {
			if err = bcrypt.CompareHashAndPassword([]byte(refreshToken.Refresh), decodeRefresh); err == nil {
				return nil
			}
		}
	}
	return err
}

func IsValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
