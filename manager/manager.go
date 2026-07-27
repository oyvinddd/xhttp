package manager

import (
	"github.com/golang-jwt/jwt/v5"
)

var (
	// access token ttl
	fifteenMinutes = 5 * time.Minute
	// refresh token ttl
	ninetyDays = 24 * 90 * time.Hour
)

type (
	Manager struct {
		secret []byte
	}
)

func NewManager(secret []byte) *Manager {
	return &Manager{secret: secret}
}

func (manager Manager) SignedAccessToken(acc account.Account) (string, error) {
	now := time.Now()
	// token expiration set to 15 minutes 
	expiration := now.Add(fifteenMinutes)
	// create custom claims (account) alongside predefined ones
	claims := account.AccessTokenClaims{
		Account: acc,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt: jwt.NewNumericDate(now),
			Subject:   acc.ID.String(),
		},
	}
	// Create a new token object, specifying signing method and the claims we would like it to contain
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(manager.secret)
}

func (manager Manager) SignedRefreshToken(accountID uuid.UUID) (string, error) {
	now := time.Now()
	// token expiration is set to 90 days
    expiration := time.Now().Add(ninetyDays)

    claims := account.RefreshTokenClaims{
        AccountID: accountID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(expiration),
            IssuedAt:  jwt.NewNumericDate(now),
            ID:        uuid.NewString(), // jti
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(manager.secret)
}

