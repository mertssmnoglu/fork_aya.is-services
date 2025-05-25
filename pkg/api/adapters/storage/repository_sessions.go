package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetSessionById(
	ctx context.Context,
	id string,
) (*users.Session, error) {
	row, err := r.queries.GetSessionById(ctx, GetSessionByIdParams{Id: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.Session{
		Id:                       row.Id,
		Status:                   row.Status,
		OauthRequestState:        row.OauthRequestState,
		OauthRequestCodeVerifier: row.OauthRequestCodeVerifier,
		OauthRedirectUri:         vars.ToStringPtr(row.OauthRedirectUri),
		LoggedInUserId:           vars.ToStringPtr(row.LoggedInUserId),
		LoggedInAt:               vars.ToTimePtr(row.LoggedInAt),
		ExpiresAt:                vars.ToTimePtr(row.ExpiresAt),
		CreatedAt:                row.CreatedAt,
		UpdatedAt:                vars.ToTimePtr(row.UpdatedAt),
	}

	return result, nil
}

func (r *Repository) CreateSession(
	ctx context.Context,
	session *users.Session,
) error {
	err := r.queries.CreateSession(ctx, CreateSessionParams{
		Id:                       session.Id,
		Status:                   session.Status,
		OauthRequestState:        session.OauthRequestState,
		OauthRequestCodeVerifier: session.OauthRequestCodeVerifier,
		OauthRedirectUri:         vars.ToSqlNullString(session.OauthRedirectUri),
		LoggedInUserId:           vars.ToSqlNullString(session.LoggedInUserId),
		LoggedInAt:               vars.ToSqlNullTime(session.LoggedInAt),
		ExpiresAt:                vars.ToSqlNullTime(session.ExpiresAt),
		CreatedAt:                session.CreatedAt,
		UpdatedAt:                vars.ToSqlNullTime(session.UpdatedAt),
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateSessionLoggedInAt(
	ctx context.Context,
	id string,
	loggedInAt time.Time,
) error {
	err := r.queries.UpdateSessionLoggedInAt(ctx, UpdateSessionLoggedInAtParams{
		Id:         id,
		LoggedInAt: sql.NullTime{Time: loggedInAt, Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}
