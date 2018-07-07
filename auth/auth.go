package auth

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http"
	"strings"
	"server/global"
)

const SecretKey = "ThisIsASecretKey"

func CreateToken(userId int64) (string, error)  {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(),
	})

	return token.SignedString([]byte(SecretKey))
}

func VerifyToken(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		tokenSlice := strings.Split(tokenString, " ")

		if len(tokenSlice) != 2 || tokenSlice[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenSlice[1], func (token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if token.Valid {
			h.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	})
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve old token from header
	tokenString := r.Header.Get("Authorization")
	tokenSlice := strings.Split(tokenString, " ")

	if len(tokenSlice) != 2 || tokenSlice[0] != "Bearer" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := jwt.Parse(tokenSlice[1], func (token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}


	userId := int64(token.Claims.(jwt.MapClaims)["userId"].(float64))

	refreshTokenString, err := CreateToken(userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json := make(map[string]interface{})
	json["token"] = refreshTokenString

	global.SendJSON(w, json, http.StatusOK)
}