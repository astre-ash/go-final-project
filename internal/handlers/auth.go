package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecretKey = "default_secret_key_for_tests"

func (h *TaskHandler) Signin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Signin error: failed to decode JSON payload: %v", err)
		sendError(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	expectedPass := os.Getenv("TODO_PASSWORD")
	if req.Password != expectedPass {
		log.Printf("Signin error: invalid password attempt")
		sendError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"hash": hashPassword(expectedPass),
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(getJWTSecret())
	if err != nil {
		log.Printf("Signin error: failed to sign JWT token: %v", err)
		sendError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "token",
		Value: signedToken,
		Path:  "/",
	})

	response := map[string]string{
		"token": signedToken,
	}
	sendJSON(w, response, http.StatusOK)
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv("TODO_PASSWORD")

		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				log.Printf("AuthMiddleware error: missing or invalid token cookie")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			token, err := jwt.Parse(cookie.Value, func(t *jwt.Token) (interface{}, error) {
				return getJWTSecret(), nil
			})

			if err != nil || !token.Valid {
				log.Printf("AuthMiddleware error: invalid token: %v", err)
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}

			tokenHash, ok := claims["hash"].(string)
			if !ok || tokenHash != hashPassword(pass) {
				log.Printf("AuthMiddleware error: failed to parse token claims")
				http.Error(w, "Authentification required", http.StatusUnauthorized)
				return
			}
		}

		next(w, r)
	})
}

func hashPassword(pass string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	return hex.EncodeToString(h.Sum(nil))
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return []byte("testSecretKey")
	}
	return []byte(secret)
}
