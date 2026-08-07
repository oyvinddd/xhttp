package auth

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

	AccessTokenClaims struct {
		// ID the account id
		ID uuid.UUID `json:"id"`
		// OrgID the organization id
		OrgID *uuid.UUID `json:"org_id"`
		// Email the account email
		Email string `json:"email"`
		// Role the account role
		Role Role `json:"role"`
		// Default claims
		jwt.RegisteredClaims
	}

	RefreshTokenClaims struct {
		// ID the account id
		ID uuid.UUID `json:"id"`
		// Default claims
		jwt.RegisteredClaims
	}
)

func GetAccessTokenClaims(ctx context.Context) (*AccessTokenClaims, bool) {
	claims, ok := ctx.Value(AccessClaimsCtxKey).(*AccessTokenClaims)
	return claims, ok
}

func (a AccessTokenClaims) HasOrganization() bool {
	return a.OrgID != nil
}

func (r Role) HasPermission(required Role) bool {
	switch r {
	case UserRole:
		return required == UserRole
	case AdminRole:
		return required == UserRole || required == AdminRole
	case GlobalRole:
		return required == UserRole || required == AdminRole || required == GlobalRole
	default:
		return false
	}
}
