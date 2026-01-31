package user

import "context"

type Client interface {
	CreateUser(ctx context.Context, in CreateUserReq) error
}
