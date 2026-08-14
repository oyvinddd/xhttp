package auth

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

	Token struct {
		// Value the actual token itself
		Value string `json:"token"`
		// ExpiresAt timestamp for token expiration 
		ExpiresAt time.Time `json:"expiration"`
	}
)

func NewManager(accessTokenTTL, refreshTokenTTL int) *Manager {
	return &Manager{accessTokenTTL, refreshTokenTTL}
}

func (m *Manager) SignedAccessToken(id uuid.UUID, orgID *uuid.UUID, email string, role Role) (Token, error) {
	now := time.Now()
	// token expiration set to 15 minutes 
	ttl := time.Duration(m.accessTokenTTL) * time.Minute
	exp := now.Add(ttl)
	// create custom claims (account) alongside predefined ones
	accessClaims := AccessTokenClaims{
		ID: id,
		OrgID: orgID,
		Email: email,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: 	jwt.NewNumericDate(exp),
			IssuedAt: 	jwt.NewNumericDate(now),
			Subject: 	id.String(),
		},
	}
	// Create a new token object, specifying signing method and the claims we would like it to contain
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	tokenStr, err := token.SignedString(hmacSecret)
	if err != nil {
		return Token{}, err
	}
	return Token{tokenStr, exp}, nil
}

func (m *Manager) SignedRefreshToken(id uuid.UUID) (Token, error) {
	now := time.Now()
	ttl := time.Duration(m.refreshTokenTTL) * 24 * time.Hour
	// token expiration is set to 90 days
    exp := now.Add(ttl)

    refreshClaims := RefreshTokenClaims{
        ID: id,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(exp),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(), // jti
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	tokenStr, err := token.SignedString(hmacSecret)
	if err != nil {
		return Token{}, err
	}
	return Token{tokenStr, exp}, nil
}

func (m *Manager) RotateRefreshToken(oldRefreshToken string) (Token, error) {
	claims, err := ParseRefreshToken(oldRefreshToken)
	if err != nil {
		return Token{}, err
	}
	return m.SignedRefreshToken(claims.ID)
}

func SetHMACSecret(secret string) {
	hmacSecret = []byte(secret)
}

func GetHMACSecret() []byte {
	return hmacSecret
}

