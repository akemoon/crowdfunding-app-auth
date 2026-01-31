package creds

import (
	"context"
	"fmt"

	"github.com/akemoon/crowdfunding-app-auth/domain"
	"github.com/akemoon/crowdfunding-app-auth/repo/creds"
	"github.com/akemoon/crowdfunding-app-auth/tool/hasher"
	"github.com/google/uuid"
)

type Service struct {
	repo   creds.Repo
	hasher hasher.Hasher
}

func NewService(repo creds.Repo, hasher hasher.Hasher) *Service {
	return &Service{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *Service) CreateCreds(ctx context.Context, req domain.SignUpRequest) (uuid.UUID, error) {
	passwordHash, err := s.hasher.Hash(req.Password)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%w: %s", domain.ErrInternal, err.Error())
	}

	userID, err := s.repo.CreateCreds(ctx, domain.Creds{
		Email:        req.Email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return uuid.Nil, fmt.Errorf("repo: %w", err)
	}

	return userID, nil
}

func (s *Service) DeleteCredsByUserID(ctx context.Context, userID uuid.UUID) error {
	err := s.repo.DeleteCredsByUserID(ctx, userID)
	if err != nil {
		return fmt.Errorf("repo: %w", err)
	}

	return nil
}

func (s *Service) ValidateCredentials(ctx context.Context, req domain.SignInRequest) (uuid.UUID, error) {
	creds, err := s.repo.GetCredsByEmail(ctx, req.Email)
	if err != nil {
		return uuid.Nil, fmt.Errorf("creds repo: %w", err)
	}

	err = s.hasher.Compare(req.Password, creds.PasswordHash)
	if err != nil {
		return uuid.Nil, fmt.Errorf("password compare: %s", err)
	}

	return creds.UserID, nil
}
