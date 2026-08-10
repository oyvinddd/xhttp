package auth

import (
	"errors"
)

var (
	ErrMissingAuthHeader 	= errors.New("missing authorization header")
	ErrInvalidAuthHeader 	= errors.New("invalid authorization header")
	ErrInvalidAccessToken 	= errors.New("invalid access token")
	ErrUnexpectedSigning 	= errors.New("unexpected signing method")
	ErrInsufficientRole 	= errors.New("insufficient role")
	ErrInvalidServiceKey 	= errors.New("invalid service key")
	ErrMissingOrganization	= errors.New("missing organization")
)

