package domain

import "github.com/google/uuid"

type CreateCredsReq struct {
	Email        string
	PasswordHash string
}

type Creds struct {
	UserID       uuid.UUID
	Email        string
	PasswordHash string
}
