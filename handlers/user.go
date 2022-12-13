package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/juanmunoz22-bit/go-rest-ws/models"
	"github.com/juanmunoz22-bit/go-rest-ws/repository"
	"github.com/juanmunoz22-bit/go-rest-ws/server"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

const (
	HASH_COST = 8
)

type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		id, err := ksuid.NewRandom()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		var user = models.User{
			Email:    request.Email,
			Password: string(hashedPassword),
			Id:       id.String(),
		}
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			http.Error(w, "User already exists", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(SignupResponse{
			Id:    user.Id,
			Email: user.Email,
		})
	}
}

func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{}
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		claim := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})
	}
}
