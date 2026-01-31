package creds

import (
	"context"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/google/uuid"
)

type Repo interface {
	CreateCreds(ctx context.Context, creds domain.Creds) (uuid.UUID, error)
	DeleteCredsByUserID(ctx context.Context, userID uuid.UUID) error
	GetCredsByEmail(ctx context.Context, email string) (domain.Creds, error)
}
