package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/maYkiss56/subscription-aggregation-service/internal/domain"
)

type SubRepository interface {
	CreateSub(ctx context.Context, sub *domain.Sub) (id uuid.UUID, err error)
	GetAllSubs(ctx context.Context) ([]*domain.Sub, error)
	GetSubByUserID(ctx context.Context, userUID uuid.UUID) ([]*domain.Sub, error)
	UpdateSub(ctx context.Context, id uuid.UUID, req *domain.UpdateSubRequest) (*domain.Sub, error)
	DeleteSub(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error)
}

type SubService struct {
	repo SubRepository
}

func New(repo SubRepository) *SubService {
	return &SubService{
		repo: repo,
	}
}

func (s *SubService) CreateSub(ctx context.Context, sub *domain.Sub) (id uuid.UUID, err error) {
	id, err = s.repo.CreateSub(ctx, sub)
	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (s *SubService) GetAllSubs(ctx context.Context) ([]*domain.Sub, error) {
	subs, err := s.repo.GetAllSubs(ctx)
	if err != nil {
		return nil, err
	}

	return subs, nil
}

func (s *SubService) GetSubByUserID(ctx context.Context, userUID uuid.UUID) ([]*domain.Sub, error) {
	sub, err := s.repo.GetSubByUserID(ctx, userUID)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubService) UpdateSub(ctx context.Context, id uuid.UUID, req *domain.UpdateSubRequest) (*domain.Sub, error) {
	sub, err := s.repo.UpdateSub(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *SubService) DeleteSub(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteSub(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s *SubService) CalculateTotalCost(ctx context.Context, filter domain.TotalCostFilter) (int, error) {
	total, err := s.repo.CalculateTotalCost(ctx, filter)
	if err != nil {
		return 0, err
	}

	return total, nil
}
