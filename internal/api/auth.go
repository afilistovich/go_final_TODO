package api

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const envPassword = "TODO_PASSWORD"

type SignInRequest struct {
	Password string `json:"password"`
}

type TokenClaims struct {
	Hash string `json:"hash"`
	jwt.RegisteredClaims
}

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pass := os.Getenv(envPassword)
		if len(pass) > 0 {
			cookie, err := r.Cookie("token")
			if err != nil {
				writeError(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			valid, err := validateJWT(cookie.Value)
			if err != nil || !valid {
				slog.Warn("Authentication failed", "error", err)
				writeError(w, "Authentication required", http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func generateJWT(password string) (string, error) {
	hash := hashPassword(password)
	claims := TokenClaims{
		Hash: hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(8 * time.Hour)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(hash))
	if err != nil {
		return "", fmt.Errorf("failed to sign jwt: %w", err)
	}
	return signedToken, nil
}

func validateJWT(tokenStr string) (bool, error) {

	password := os.Getenv(envPassword)
	hash := hashPassword(password)
	claims := &TokenClaims{}

	jwtToken, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(hash), nil
	}, jwt.WithExpirationRequired(),
	)

	if err != nil {
		slog.Warn("Failed to parse JWT", "error", err)
		return false, err
	}

	if !jwtToken.Valid {
		slog.Warn("Invalid JWT token")
		return false, nil
	}

	if claims.Hash != hash {
		slog.Warn("JWT hash mismatch (password changed?)")
		return false, nil
	}
	return true, nil
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	var req SignInRequest
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Invalid JSON in request", "error", err)
		writeError(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	password := os.Getenv(envPassword)
	if req.Password != password {
		writeError(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(req.Password)
	if err != nil {
		slog.Error("Failed to create jwt", "error", err)
		writeError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err = writeJSON(w, http.StatusOK, map[string]string{"token": token}); err != nil {
		slog.Error("Failed to encode token response", "error", err)
		return
	}
}
