package midleware

import (
	"backend/internal/auth"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

type ctxKey string

const CtxUserIDKey ctxKey = "userID"

func JWTMidleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]
		token, err := auth.ParseToken(tokenStr)
		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		sub, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "invalid subject claim", http.StatusUnauthorized)
			return
		}
		userID := int(sub)
		ctx := context.WithValue(r.Context(), CtxUserIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
