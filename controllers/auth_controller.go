package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"qisur-challenge/config"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Solicitud inválida", http.StatusBadRequest)
			return
		}

		if creds.Username != "admin" || creds.Password != "admin" {
			http.Error(w, "Credenciales inválidas", http.StatusUnauthorized)
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": creds.Username,
			"exp":      time.Now().Add(1 * time.Hour).Unix(),
		})

		tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
		if err != nil {
			http.Error(w, "Error al generar token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
		})
	}
}
