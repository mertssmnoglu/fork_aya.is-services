package custom_domains

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/ajan/logfx"
)

var (
	ErrFailedToGetRecord = errors.New("failed to get record")
	// ErrFailedToListRecords = errors.New("failed to list records")
	// ErrFailedToCreateRecord = errors.New("failed to create record").
)

type Repository interface {
	GetCustomDomainByDomain(ctx context.Context, domain string) (*GetCustomDomainByDomainRow, error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetByDomain(ctx context.Context, domain string) (*GetCustomDomainByDomainRow, error) {
	record, err := s.repo.GetCustomDomainByDomain(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("%w(domain: %s): %w", ErrFailedToGetRecord, domain, err)
	}

	return record, nil
}
