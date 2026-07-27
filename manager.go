package xtoken

import (
	"strings"
	"net/http"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

type (
	Manager struct {
		secret []byte
		accessTokenTTL int
		refreshTokenTTL int
	}
)

func New(secret []byte, accessTokenTTL, refreshTokenTTL int) *Manager {
	return &Manager{secret, accessTokenTTL, refreshTokenTTL}
}

func (manager Manager) SignedAccessToken(id uuid.UUID, email string, role Role) (string, error) {
	now := time.Now()
	// token expiration set to 15 minutes 
	expiration := now.Add(manager.accessTokenTTL)
	// create custom claims (account) alongside predefined ones
	accessClaims := AccessTokenClaims{
		ID: id.String(),
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: 	jwt.NewNumericDate(expiration),
			IssuedAt: 	jwt.NewNumericDate(now),
			Subject: 	id.String(),
		},
	}
	// Create a new token object, specifying signing method and the claims we would like it to contain
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    return token.SignedString(manager.secret)
}

func (manager Manager) SignedRefreshToken(id uuid.UUID) (string, error) {
	now := time.Now()
	// token expiration is set to 90 days
    expiration := time.Now().Add(manager.refreshTokenTTL)

    refreshClaims := RefreshTokenClaims{
        ID: id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiration),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(), // jti
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
    return token.SignedString(manager.secret)
}

func (manager Manager) Authorize(next httprouter.Handle, role Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenStr, err := parseAccessToken(r.Header.Get("Authorization"))
		if err != nil {
			http.Error(w, err, http.StatusUnauthorized)
			return
		}

		accessClaims := new(AccessClaims)

		token, err := jwt.ParseWithClaims(tokenStr, accessClaims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, unexpectedSigningError
			}
			return manager.secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, invalidAccessTokenError, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), claims.AccessClaimsCtxKey, accessClaims)
		next(w, r.WithContext(ctx), ps)
	}
}

func parseAccessToken(header string) (string, error) {
	if header == "" {
		return "", missingAuthHeaderError
	}
	prefix := "Bearer "
	if !strings.HasPrefix(header, prefix) {
		return "", invalidAuthHeaderError
	}
	return strings.TrimPrefix(header, prefix)
}

