package auth

import (
	"fmt"
	"os"
	"time"

	"backend/internal/config"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateToken(userID int) (string, error) {
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed loading config: %w", err)
	}
	privData, err := os.ReadFile(cfg.JWTPrivateKeyPath)
	if err != nil {
		return "", fmt.Errorf("failed reading private key: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privData)
	if err != nil {
		return "", fmt.Errorf("failed parsing private key: %w", err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"sub": userID,
		"iss": "horoscope-vault",
		"aud": "vault-users",
		"exp": jwt.NewNumericDate(jwt.TimeFunc().Add(24 * time.Hour)),
	})
	tokenStr, err := token.SignedString(privKey)
	if err != nil {
		return "", fmt.Errorf("failed signing token: %w", err)
	}
	return tokenStr, nil
}

func ParseToken(tokenStr string) (jwt.Token, error) {
	cfg, err := config.Load()
	if err != nil {
		return jwt.Token{}, fmt.Errorf("failed loading config: %w", err)
	}
	pubData, err := os.ReadFile(cfg.JWTPublicKeyPath)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("failed reading public key: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubData)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("failed parsing public key: %w", err)
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return jwt.Token{}, fmt.Errorf("token parse error: %w", err)
	}
	return *token, nil
}
