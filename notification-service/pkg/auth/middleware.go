package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

type contextKey string

const userIDKey contextKey = "userID"

func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(userIDKey)
	id, ok := val.(uuid.UUID)
	return id, ok
}

func SetUserIDInContext(ctx context.Context, userID uuid.UUID) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func Middleware(jwtManager *JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			http.Error(w, ErrInvalidAuthHeader.Error(), http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, bearerPrefix)

		if token == "" {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		data, err := jwtManager.Parse(token)
		if err != nil {
			http.Error(w, ErrInvalidToken.Error(), http.StatusUnauthorized)
			return
		}

		ctx := SetUserIDInContext(r.Context(), data.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
