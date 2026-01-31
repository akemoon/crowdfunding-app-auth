package handler

import (
	"errors"
	"net/http"

	"github.com/akemoon/crowdfunding-app-auth/cluster/user"
	"github.com/akemoon/crowdfunding-app-auth/domain"
)

const (
	HttpErrEmailExists        = "email_exists"
	HttpErrUsernameExists     = "username_exists"
	HttpErrUnknownConflict    = "unknown_conflict"
	HttpInternalError         = "internal_error"
	HttpErrInvalidAccessToken = "invalid_access_token"
)

type ErrResp struct {
	Error   string `json:"error"`
	Details string `json:"details"`
}

func mapErrToHTTP(err error) (int, ErrResp) {

	if errors.Is(err, domain.ErrEmailExists) {
		return http.StatusConflict, ErrResp{
			Error:   HttpErrEmailExists,
			Details: domain.ErrEmailExists.Error(),
		}
	}

	if errors.Is(err, user.ErrUsernameExists) {
		return http.StatusConflict, ErrResp{
			Error:   HttpErrUsernameExists,
			Details: user.ErrUsernameExists.Error(),
		}
	}

	if errors.Is(err, domain.ErrUnknownConflict) {
		return http.StatusConflict, ErrResp{
			Error:   HttpErrUnknownConflict,
			Details: domain.ErrUnknownConflict.Error(),
		}
	}

	if errors.Is(err, domain.ErrInvalidAccessToken) {
		return http.StatusUnauthorized, ErrResp{
			Error:   HttpErrInvalidAccessToken,
			Details: domain.ErrInvalidAccessToken.Error(),
		}
	}

	return http.StatusInternalServerError, ErrResp{
		Error:   HttpInternalError,
		Details: domain.ErrInternal.Error(),
	}
}
