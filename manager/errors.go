package manager

import (
	"errors"
)

var (
	missingAuthHeaderError error = errors.New("missing authorization header")
	invalidAuthHeaderError error = errors.New("invalid authorization header")
	invalidAccessTokenError error = errors.New("invalid access token")
	unexpectedSigningError error = errors.New("unexpected signing method")
)

