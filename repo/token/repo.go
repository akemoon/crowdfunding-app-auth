package token

import (
	"context"
)

type RefreshTokenRepo interface {
	Set(ctx context.Context, refreshToken string) error
	Check(ctx context.Context, refreshToken string) error
	Delete(ctx context.Context, refreshToken string) error
}
