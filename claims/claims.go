package claims

import (
	"context"
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v5"
)

const (
	UserRole Role = 0

	AdminRole Role = 1

	GlobalRole Role = 2
)

const (
	AccessClaimsCtxKey string = "accessClaims"
)

type (
	Role int

	AccessClaims struct {
		// ID the account id
		ID uuid.UUID `json:"id"`
		// Email the account email
		Email string `json:"email"`
		// Role the account role
		Role Role `json:"role"`
		// Default claims
		jwt.RegisteredClaims
	}

	RefreshClaims struct {
		// ID the account id
		ID uuid.UUID `json:"id"`
		// Default claims
		jwt.RegisteredClaims
	}
)

func GetAccessClaims(ctx context.Context) (*AccessClaims, bool) {
	claims, ok := ctx.Value(ClaimsContextKey).(*AccessClaims)
	return claims, ok
}
