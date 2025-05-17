package storage

import (
	"context"
	"database/sql"

	"github.com/eser/aya.is-services/pkg/api/business/users"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetUserById(
	ctx context.Context,
	localeCode string,
	id string,
) (*users.User, error) {
	row, err := r.queries.GetUserById(ctx, GetUserByIdParams{Id: id})
	if err != nil {
		return nil, err
	}

	result := &users.User{
		Id:                  row.Id,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
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
	localeCode string,
	email string,
) (*users.User, error) {
	row, err := r.queries.GetUserByEmail(
		ctx,
		GetUserByEmailParams{Email: sql.NullString{String: email, Valid: true}},
	)
	if err != nil {
		return nil, err
	}

	result := &users.User{
		Id:                  row.Id,
		Kind:                row.Kind,
		Name:                row.Name,
		Email:               vars.ToStringPtr(row.Email),
		Phone:               vars.ToStringPtr(row.Phone),
		GithubHandle:        vars.ToStringPtr(row.GithubHandle),
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
	localeCode string,
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
