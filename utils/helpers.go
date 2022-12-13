package utils

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/juanmunoz22-bit/go-rest-ws/models"
	"github.com/juanmunoz22-bit/go-rest-ws/server"
)

func StringToken(r *http.Request) (tokenString string) {
	return strings.TrimSpace(r.Header.Get("Authorization"))
}

func ClaimTokenFromAuth(s server.Server, tokenString string, claim models.AppClaims) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.Config().JWTSecret), nil
	})
}
