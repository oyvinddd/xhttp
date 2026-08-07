package auth

import (
	"errors"
)

var (
	MissingAuthHeaderError = errors.New("missing authorization header")
	InvalidAuthHeaderError = errors.New("invalid authorization header")
	InvalidAccessTokenError = errors.New("invalid access token")
	UnexpectedSigningError = errors.New("unexpected signing method")
	InsufficientRoleError = errors.New("insufficient role")
	InvalidServiceKeyError = errors.New("invalid service key")
)

