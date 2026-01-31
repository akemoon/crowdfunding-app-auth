package token

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/akemoon/crowdfunding-app-auth/repo/token"
	"github.com/golang-jwt/jwt/v5"
)

const (
	accessTokenLifeTime  = 15 * time.Minute
	refreshTokenLifeTime = 7 * 24 * time.Hour
)

type Service struct {
	refreshTokenRepo token.RefreshTokenRepo
	secret           string
}

func NewService(r token.RefreshTokenRepo, s string) *Service {
	return &Service{
		refreshTokenRepo: r,
		secret:           s,
	}
}

func (s *Service) GenAccessToken(ctx context.Context, tc domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": tc.UserID,
		"exp":    time.Now().Add(accessTokenLifeTime).Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	return signedToken, nil
}

func (s *Service) GenRefreshToken(ctx context.Context, tc domain.TokenClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": tc.UserID,
		"exp":    time.Now().Add(refreshTokenLifeTime).Unix(),
	})

	signedToken, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}

	err = s.refreshTokenRepo.Set(ctx, signedToken)
	if err != nil {
		return "", fmt.Errorf("token repo: %w", err)
	}

	return signedToken, nil
}

//func (s *Service) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
//	err := s.tokenRepo.Check(ctx, refreshToken)
//	if err != nil {
//		return "", "", err
//	}
//
//	userID, err := s.userIDFromRefreshToken(refreshToken)
//	if err != nil {
//		return "", "", err
//	}
//
//	at, err := s.GenAccessToken(ctx, domain.TokenClaims{
//		UserID: userID,
//	})
//	if err != nil {
//		return "", "", err
//	}
//
//	rt, err := s.GenRefreshToken(ctx, domain.TokenClaims{
//		UserID: userID,
//	})
//	if err != nil {
//		return "", "", err
//	}
//
//	return at, rt, nil
//}
//
//func (s *Service) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
//	return s.tokenRepo.Delete(ctx, refreshToken)
//}
//
//func (s *Service) DeleteUserTokens(ctx context.Context, refreshToken string) error {
//	userID, err := s.userIDFromRefreshToken(refreshToken)
//	if err != nil {
//		return err
//	}
//
//	return s.tokenRepo.DeleteByUserID(ctx, userID)
//}

func (s *Service) ValidateAccessToken(token string) error {
	parts := strings.Fields(token)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return domain.ErrInvalidAccessToken
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	jwtToken, err := parser.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("%w: unexpected signing method: %s", domain.ErrInvalidAccessToken, t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return fmt.Errorf("%w: %s", domain.ErrInternal, err)
	}
	if !jwtToken.Valid {
		return domain.ErrInvalidAccessToken
	}

	return nil
}
