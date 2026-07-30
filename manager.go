package xtoken

import (
	"time"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

var (
	hmacSecret []byte 
)

type (
	Manager struct {
		accessTokenTTL int
		refreshTokenTTL int
	}
)

func NewManager(accessTokenTTL, refreshTokenTTL int) *Manager {
	return &Manager{accessTokenTTL, refreshTokenTTL}
}

func (manager Manager) SignedAccessToken(id uuid.UUID, email string, role Role) (string, error) {
	now := time.Now()
	// token expiration set to 15 minutes 
	ttl := time.Duration(manager.accessTokenTTL) * time.Minute
	expiration := now.Add(ttl)
	// create custom claims (account) alongside predefined ones
	accessClaims := AccessTokenClaims{
		ID: id,
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
    return token.SignedString(hmacSecret)
}

func (manager Manager) SignedRefreshToken(id uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	ttl := time.Duration(manager.refreshTokenTTL) * 24 * time.Hour
	// token expiration is set to 90 days
    expiration := time.Now().Add(ttl)

    refreshClaims := RefreshTokenClaims{
        ID: id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiration),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(), // jti
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedTokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", time.Now(), err
	}
	return signedTokenString, expiration, nil
}

func SetHMACSecret(secret string) {
	hmacSecret = []byte(secret)
}

func GetHMACSecret() []byte {
	return hmacSecret
}

