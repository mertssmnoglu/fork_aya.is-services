package stories

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord   = errors.New("failed to get record")
	ErrFailedToListRecords = errors.New("failed to list records")
	// ErrFailedToCreateRecord = errors.New("failed to create record").
)

type Repository interface {
	GetProfileIdBySlug(ctx context.Context, slug string) (string, error)
	GetProfileById(ctx context.Context, localeCode string, id string) (*profiles.Profile, error)
	GetStoryIdBySlug(ctx context.Context, slug string) (string, error)
	GetStoryById(ctx context.Context, localeCode string, id string) (*Story, error)
	GetStoryBySlug(ctx context.Context, localeCode string, slug string) (*Story, error)
	ListStories(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Story], error)
	ListStoriesByAuthorProfileId(
		ctx context.Context,
		localeCode string,
		authorProfileId string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Story], error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetById(ctx context.Context, localeCode string, id string) (*Story, error) {
	record, err := s.repo.GetStoryById(ctx, localeCode, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetBySlug(ctx context.Context, localeCode string, slug string) (*Story, error) {
	record, err := s.repo.GetStoryBySlug(ctx, localeCode, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return record, nil
}

func (s *Service) GetBySlugEx(
	ctx context.Context,
	localeCode string,
	slug string,
) (*StoryWithChildren, error) {
	storyId, err := s.repo.GetStoryIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetStoryById(ctx, localeCode, storyId)
	if err != nil {
		return nil, fmt.Errorf("%w(story_id: %s): %w", ErrFailedToGetRecord, storyId, err)
	}

	authorProfileId := *record.AuthorProfileId

	authorProfile, err := s.repo.GetProfileById(ctx, localeCode, authorProfileId)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, authorProfileId, err)
	}

	result := &StoryWithChildren{
		Story:         record,
		AuthorProfile: authorProfile,
	}

	return result, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Story], error) {
	records, err := s.repo.ListStories(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*Story]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) ListByAuthorProfileSlug(
	ctx context.Context,
	localeCode string,
	authorProfileSlug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Story], error) {
	authorProfileId, err := s.repo.GetProfileIdBySlug(ctx, authorProfileSlug)
	if err != nil {
		return cursors.Cursored[[]*Story]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			authorProfileSlug,
			err,
		)
	}

	records, err := s.repo.ListStoriesByAuthorProfileId(
		ctx,
		localeCode,
		authorProfileId,
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*Story]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}
