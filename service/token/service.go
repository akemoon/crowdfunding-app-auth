package token

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/akemoon/crowdfunding-app-auth/repo/token"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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
	// TODO: user id : token or just token
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

func (s *Service) DeleteRefreshToken(ctx context.Context, refreshToken string) error {
	err := s.refreshTokenRepo.Delete(ctx, refreshToken)
	if err != nil {
		return fmt.Errorf("token repo: %w", err)
	}

	return nil
}

func (s *Service) ValidateAccessToken(token string) (uuid.UUID, error) {
	parts := strings.Fields(token)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return uuid.Nil, domain.ErrInvalidAccessToken
	}

	parser := jwt.NewParser(jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	jwtToken, err := parser.Parse(parts[1], func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("%w: unexpected signing method: %s", domain.ErrInvalidAccessToken, t.Header["alg"])
		}
		return []byte(s.secret), nil
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("parse token: %w", err)
	}
	if !jwtToken.Valid {
		return uuid.Nil, domain.ErrInvalidAccessToken
	}

	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, domain.ErrInvalidAccessToken
	}

	userIDStr, ok := claims["userID"].(string)
	if !ok || userIDStr == "" {
		return uuid.Nil, domain.ErrInvalidAccessToken
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, domain.ErrInvalidAccessToken
	}

	return userID, nil
}
