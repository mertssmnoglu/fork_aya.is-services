package profiles

import (
	"context"
	"errors"
	"fmt"

	"github.com/eser/ajan/logfx"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
)

var (
	ErrFailedToGetRecord   = errors.New("failed to get record")
	ErrFailedToListRecords = errors.New("failed to list records")
	// ErrFailedToCreateRecord = errors.New("failed to create record").
)

type RecentPostsFetcher interface {
	GetRecentPostsByUsername(
		ctx context.Context,
		username string,
		userId string,
	) ([]*ExternalPost, error)
}

type Repository interface {
	GetProfileIdBySlug(ctx context.Context, slug string) (string, error)
	GetProfileIdByCustomDomain(ctx context.Context, domain string) (string, error)
	GetProfileById(ctx context.Context, localeCode string, id string) (*Profile, error)
	ListProfiles(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*Profile], error)
	// GetProfileLinksForKind(ctx context.Context, kind string) ([]*ProfileLink, error)
	GetProfilePagesByProfileId(
		ctx context.Context,
		localeCode string,
		profileId string,
	) ([]*ProfilePageBrief, error)
	GetProfilePageByProfileIdAndSlug(
		ctx context.Context,
		localeCode string,
		profileId string,
		pageSlug string,
	) (*ProfilePage, error)
	GetProfileLinksByProfileId(
		ctx context.Context,
		localeCode string,
		profileId string,
	) ([]*ProfileLinkBrief, error)
	ListProfileContributions(
		ctx context.Context,
		localeCode string,
		profileId string,
		kinds []string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*ProfileMembership], error)
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
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetProfileById(ctx, localeCode, profileId)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	return record, nil
}

func (s *Service) GetBySlugEx(
	ctx context.Context,
	localeCode string,
	slug string,
) (*ProfileWithChildren, error) {
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	record, err := s.repo.GetProfileById(ctx, localeCode, profileId)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	pages, err := s.repo.GetProfilePagesByProfileId(ctx, localeCode, record.Id)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	links, err := s.repo.GetProfileLinksByProfileId(ctx, localeCode, record.Id)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	result := &ProfileWithChildren{
		Profile: record,
		Pages:   pages,
		Links:   links,
	}

	return result, nil
}

func (s *Service) GetByCustomDomain(
	ctx context.Context,
	localeCode string,
	domain string,
) (*Profile, error) {
	profileId, err := s.repo.GetProfileIdByCustomDomain(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("%w(custom_domain: %s): %w", ErrFailedToGetRecord, domain, err)
	}

	record, err := s.repo.GetProfileById(ctx, localeCode, profileId)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*Profile], error) {
	records, err := s.repo.ListProfiles(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*Profile]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}

func (s *Service) ListPagesBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) ([]*ProfilePageBrief, error) {
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	pages, err := s.repo.GetProfilePagesByProfileId(ctx, localeCode, profileId)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	return pages, nil
}

func (s *Service) GetPageBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	pageSlug string,
) (*ProfilePage, error) {
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	page, err := s.repo.GetProfilePageByProfileIdAndSlug(ctx, localeCode, profileId, pageSlug)
	if err != nil {
		return nil, fmt.Errorf(
			"%w(profile_id: %s, page_slug: %s): %w",
			ErrFailedToGetRecord,
			profileId,
			pageSlug,
			err,
		)
	}

	return page, nil
}

func (s *Service) ListLinksBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
) ([]*ProfileLinkBrief, error) {
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("%w(slug: %s): %w", ErrFailedToGetRecord, slug, err)
	}

	links, err := s.repo.GetProfileLinksByProfileId(ctx, localeCode, profileId)
	if err != nil {
		return nil, fmt.Errorf("%w(profile_id: %s): %w", ErrFailedToGetRecord, profileId, err)
	}

	return links, nil
}

func (s *Service) ListProfileContributionsBySlug(
	ctx context.Context,
	localeCode string,
	slug string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*ProfileMembership], error) {
	profileId, err := s.repo.GetProfileIdBySlug(ctx, slug)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w(slug: %s): %w",
			ErrFailedToGetRecord,
			slug,
			err,
		)
	}

	memberships, err := s.repo.ListProfileContributions(
		ctx,
		localeCode,
		profileId,
		[]string{"organization", "product"},
		cursor,
	)
	if err != nil {
		return cursors.Cursored[[]*ProfileMembership]{}, fmt.Errorf(
			"%w: %w",
			ErrFailedToListRecords,
			err,
		)
	}

	return memberships, nil
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
