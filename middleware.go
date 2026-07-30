package xtoken

import (
	"net/http"
	"strings"
	"context"
	"github.com/golang-jwt/jwt/v5"
)

func (manager Manager) Authorize(next http.Handler, requiredRole Role) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr, err := ParseAccessToken(r.Header.Get("Authorization"))
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

