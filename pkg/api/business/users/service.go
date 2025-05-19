package users

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

type Repository interface {
	GetUserById(ctx context.Context, localeCode string, id string) (*User, error)
	GetUserByEmail(ctx context.Context, localeCode string, email string) (*User, error)
	ListUsers(
		ctx context.Context,
		localeCode string,
		cursor *cursors.Cursor,
	) (cursors.Cursored[[]*User], error)
}

// --- Auth & OAuth Ports ---

type OAuthService interface {
	// InitiateOAuth returns the URL to redirect the user to GitHub for login, and the state to track the request.
	InitiateOAuth(ctx context.Context, redirectURI string) (authURL string, state OAuthState, err error)

	// HandleOAuthCallback exchanges the code for a token, fetches user info, upserts user, creates session, and returns JWT.
	HandleOAuthCallback(ctx context.Context, code string, state string) (AuthResult, error)
}

type SessionService interface {
	// UpdateLoggedInAt updates the session's logged_in_at timestamp.
	UpdateLoggedInAt(ctx context.Context, sessionID string) error

	// ValidateJWT parses and validates the JWT, returning claims if valid.
	ValidateJWT(token string) (JWTClaims, error)
}

type Service struct {
	logger      *logfx.Logger
	repo        Repository
	idGenerator RecordIDGenerator
}

func NewService(logger *logfx.Logger, repo Repository) *Service {
	return &Service{logger: logger, repo: repo, idGenerator: DefaultIDGenerator}
}

func (s *Service) GetById(ctx context.Context, localeCode string, id string) (*User, error) {
	record, err := s.repo.GetUserById(ctx, localeCode, id)
	if err != nil {
		return nil, fmt.Errorf("%w(id: %s): %w", ErrFailedToGetRecord, id, err)
	}

	return record, nil
}

func (s *Service) GetByEmail(ctx context.Context, localeCode string, email string) (*User, error) {
	record, err := s.repo.GetUserByEmail(ctx, localeCode, email)
	if err != nil {
		return nil, fmt.Errorf("%w(email: %s): %w", ErrFailedToGetRecord, email, err)
	}

	return record, nil
}

func (s *Service) List(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*User], error) {
	records, err := s.repo.ListUsers(ctx, localeCode, cursor)
	if err != nil {
		return cursors.Cursored[[]*User]{}, fmt.Errorf("%w: %w", ErrFailedToListRecords, err)
	}

	return records, nil
}
