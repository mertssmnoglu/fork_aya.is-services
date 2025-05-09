package profiles

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

type RecentPostsFetcher interface {
	GetRecentPostsByUsername(ctx context.Context, username string, userId string) ([]*ExternalPost, error)
}

type Repository interface {
	GetProfileById(ctx context.Context, localeCode string, id string) (*Profile, error)
	GetProfileBySlug(ctx context.Context, localeCode string, slug string) (*Profile, error)
	GetProfileByCustomDomain(ctx context.Context, localeCode string, domain string) (*Profile, error)
	ListProfiles(ctx context.Context, localeCode string) ([]*Profile, error)
	// GetProfileLinksForKind(ctx context.Context, kind string) ([]*ProfileLink, error)
	// CreateProfile(ctx context.Context, arg CreateProfileParams) (*Profile, error)
	// UpdateProfile(ctx context.Context, arg UpdateProfileParams) (int64, error)
	// DeleteProfile(ctx context.Context, id string) (int64, error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetById(ctx context.Context, localeCode string, id string) (*Profile, error) {
	record, err := s.repo.GetProfileById(ctx, localeCode, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(ctx context.Context, localeCode string, slug string) (*Profile, error) {
	record, err := s.repo.GetProfileBySlug(ctx, localeCode, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return record, nil
}

func (s *Service) GetByCustomDomain(ctx context.Context, localeCode string, domain string) (*Profile, error) {
	record, err := s.repo.GetProfileByCustomDomain(ctx, localeCode, domain)
	if err != nil {
		return nil, fmt.Errorf("%w(custom_domain: %s): %w", ErrFailedToGetRecord, domain, err)
	}

	return record, nil
}

func (s *Service) List(ctx context.Context, localeCode string) ([]*Profile, error) {
	records, err := s.repo.ListProfiles(ctx, localeCode)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) Import(ctx context.Context, fetcher RecentPostsFetcher) error {
	// 	links, err := s.repo.GetProfileLinksForKind(ctx, "x")
	// 	if err != nil {
	// 		return fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	// 	}
	// 	for _, link := range links {
	// 		s.logger.InfoContext(ctx, "importing posts", "kind", link.Kind, "title", link.Title)
	// 		posts, err := fetcher.GetRecentPostsByUsername(ctx, link.RemoteId.String, link.AuthAccessToken)
	// 		if err != nil {
	// 			return fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	// 		}
	// 		s.logger.InfoContext(ctx, "posts imported", "kind", link.Kind, "title", link.Title, "posts", posts)
	// 	}
	return nil
}

// func (s *Service) Create(ctx context.Context, input *Profile) (*Profile, error) {
// 	record, err := s.repo.CreateProfile(ctx, input)
// 	if err != nil {
// 		return nil, fmt.Errorf("%w: %w", ErrFailedToCreateRecord, err)
// 	}

// 	return record, nil
// }
