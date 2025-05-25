package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetUserById(
	ctx context.Context,
	id string,
) (*users.User, error) {
	row, err := r.queries.GetUserById(ctx, GetUserByIdParams{Id: id})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.User{
		Id:                  row.Id,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
		GithubRemoteId:      vars.ToStringPtr(row.GithubRemoteId),
		BskyHandle:          vars.ToStringPtr(row.BskyHandle),
		XHandle:             vars.ToStringPtr(row.XHandle),
		IndividualProfileId: vars.ToStringPtr(row.IndividualProfileId),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetUserByEmail(
	ctx context.Context,
	email string,
) (*users.User, error) {
	row, err := r.queries.GetUserByEmail(
		ctx,
		GetUserByEmailParams{Email: sql.NullString{String: email, Valid: true}},
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &users.User{
		Id:                  row.Id,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
		GithubRemoteId:      vars.ToStringPtr(row.GithubRemoteId),
		BskyHandle:          vars.ToStringPtr(row.BskyHandle),
		XHandle:             vars.ToStringPtr(row.XHandle),
		IndividualProfileId: vars.ToStringPtr(row.IndividualProfileId),
		CreatedAt:           row.CreatedAt,
		UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
		DeletedAt:           vars.ToTimePtr(row.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListUsers(
	ctx context.Context,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*users.User], error) {
	var wrappedResponse cursors.Cursored[[]*users.User]

	rows, err := r.queries.ListUsers(
		ctx,
		ListUsersParams{
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*users.User, len(rows))
	for i, row := range rows {
		result[i] = &users.User{
			Id:                  row.Id,
			Kind:                row.Kind,
			Name:                row.Name,
			Email:               vars.ToStringPtr(row.Email),
			Phone:               vars.ToStringPtr(row.Phone),
			GithubHandle:        vars.ToStringPtr(row.GithubHandle),
			GithubRemoteId:      vars.ToStringPtr(row.GithubRemoteId),
			BskyHandle:          vars.ToStringPtr(row.BskyHandle),
			XHandle:             vars.ToStringPtr(row.XHandle),
			IndividualProfileId: vars.ToStringPtr(row.IndividualProfileId),
			CreatedAt:           row.CreatedAt,
			UpdatedAt:           vars.ToTimePtr(row.UpdatedAt),
			DeletedAt:           vars.ToTimePtr(row.DeletedAt),
		}
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].Id
	}

	return wrappedResponse, nil
}

func (r *Repository) CreateUser(
	ctx context.Context,
	user *users.User,
) error {
	err := r.queries.CreateUser(ctx, CreateUserParams{
		Id:                  user.Id,
		Kind:                user.Kind,
		Name:                user.Name,
		Email:               vars.ToSqlNullString(user.Email),
		Phone:               vars.ToSqlNullString(user.Phone),
		GithubHandle:        vars.ToSqlNullString(user.GithubHandle),
		GithubRemoteId:      vars.ToSqlNullString(user.GithubRemoteId),
		BskyHandle:          vars.ToSqlNullString(user.BskyHandle),
		BskyRemoteId:        sql.NullString{String: "", Valid: false},
		XHandle:             vars.ToSqlNullString(user.XHandle),
		XRemoteId:           sql.NullString{String: "", Valid: false},
		IndividualProfileId: vars.ToSqlNullString(user.IndividualProfileId),
		CreatedAt:           user.CreatedAt,
		UpdatedAt:           vars.ToSqlNullTime(user.UpdatedAt),
		DeletedAt:           vars.ToSqlNullTime(user.DeletedAt),
	})
	if err != nil {
		return err
	}

	return nil
}
