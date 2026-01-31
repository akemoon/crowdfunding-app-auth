package auth

import (
	"context"
	"fmt"

	"github.com/akemoon/crowdfunding-app-auth/cluster/user"
	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/akemoon/crowdfunding-app-auth/service/creds"
	"github.com/akemoon/crowdfunding-app-auth/service/token"
)

type Service struct {
	userClient user.Client
	credsSvc   *creds.Service
	tokenSvc   *token.Service
}

func NewService(uc user.Client, cs *creds.Service, ts *token.Service) *Service {
	return &Service{
		userClient: uc,
		credsSvc:   cs,
		tokenSvc:   ts,
	}
}

func (s *Service) SignUp(ctx context.Context, req domain.SignUpRequest) error {
	userID, err := s.credsSvc.CreateCreds(ctx, req)
	if err != nil {
		return fmt.Errorf("creds service: %w", err)
	}

	err = s.userClient.CreateUser(ctx, user.CreateUserReq{
		UserID:   userID,
		Username: req.Username,
	})
	if err != nil {
		rollbackErr := s.credsSvc.DeleteCredsByUserID(ctx, userID)
		if rollbackErr != nil {
			return fmt.Errorf("creds service: %w", err)
		}

		return fmt.Errorf("user client: %w", err)
	}

	return nil
}

func (s *Service) SignIn(ctx context.Context, req domain.SignInRequest) (domain.SignInResponse, error) {
	userID, err := s.credsSvc.ValidateCredentials(ctx, req)
	if err != nil {
		return domain.SignInResponse{}, fmt.Errorf("creds service: %w", err)
	}

	accessToken, err := s.tokenSvc.GenAccessToken(ctx, domain.TokenClaims{
		UserID: userID,
	})
	if err != nil {
		return domain.SignInResponse{}, fmt.Errorf("token service: %w", err)
	}

	refreshToken, err := s.tokenSvc.GenRefreshToken(ctx, domain.TokenClaims{
		UserID: userID,
	})
	if err != nil {
		return domain.SignInResponse{}, fmt.Errorf("token service: %w", err)
	}

	return domain.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// func (s *Service) Logout(ctx context.Context, refreshToken string) error {
// 	if refreshToken == "" {
// 		return domain.ErrInvalidCreds
// 	}

// 	if err := s.tokenSvc.DeleteUserTokens(ctx, refreshToken); err != nil {
// 		return err
// 	}

// 	return nil
// }
