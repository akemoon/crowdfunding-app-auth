package user

import "github.com/google/uuid"

type CreateUserReq struct {
	UserID   uuid.UUID `json:"userID"`
	Username string    `json:"username"`
}

type CreateUserResp struct {
	UserID uuid.UUID `json:"userID"`
}
