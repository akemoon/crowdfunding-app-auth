package domain

import "errors"

var (
	ErrEmailExists         = errors.New("email already exists")
	ErrUnknownConflict     = errors.New("unknown conflict")
	ErrCredsNotFound       = errors.New("creds not found")
	ErrInvalidPassrord     = errors.New("invalid password")
	ErrInvalidAccessToken  = errors.New("invalid access token")
	ErrInvlaidRefreshToken = errors.New("invalid refresh token")

	ErrInternal = errors.New("internal error")
)
