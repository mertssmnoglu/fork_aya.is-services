package stories

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/ajan/logfx"
)

var (
	ErrFailedToGetRecord   = errors.New("failed to get record")
	ErrFailedToListRecords = errors.New("failed to list records")
	// ErrFailedToCreateRecord = errors.New("failed to create record").
)

type Repository interface {
	GetStoryById(ctx context.Context, id string) (*Story, error)
	GetStoryBySlug(ctx context.Context, slug string) (*Story, error)
	ListStories(ctx context.Context) ([]*Story, error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetById(ctx context.Context, id string) (*Story, error) {
	record, err := s.repo.GetStoryById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(ctx context.Context, slug string) (*Story, error) {
	record, err := s.repo.GetStoryBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return record, nil
}

func (s *Service) List(ctx context.Context) ([]*Story, error) {
	records, err := s.repo.ListStories(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}
