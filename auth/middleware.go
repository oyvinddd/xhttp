package auth

import (
	"net/http"
	"strings"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

const (
	authorizationHeader = "Authorization"
	serviceKeyHeader = "X-Service-Key"
)

func Authorize(next http.Handler, requiredRole Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := ParseAccessToken(r.Header.Get(authorizationHeader))
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		accessClaims := new(AccessTokenClaims)

		token, err := jwt.ParseWithClaims(tokenStr, accessClaims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, UnexpectedSigningError
			}
			return hmacSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, InvalidAccessTokenError.Error(), http.StatusUnauthorized)
			return
		}

		if !accessClaims.Role.HasPermission(requiredRole) {
			http.Error(w, InsufficientRoleError.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AccessClaimsCtxKey, accessClaims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequireServiceKey(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(serviceKeyHeader) != key {
				http.Error(w, InvalidServiceKeyError.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ParseAccessToken(header string) (string, error) {
	if header == "" {
		return "", MissingAuthHeaderError
	}
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", InvalidAuthHeaderError
	}
	return strings.TrimPrefix(header, prefix), nil
}

