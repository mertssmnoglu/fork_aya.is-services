//nolint:dupl
package storage

import (
	"context"
	"database/sql"

	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetProfileById(ctx context.Context, localeCode string, id string) (*profiles.Profile, error) {
	row, err := r.queries.GetProfileById(ctx, GetProfileByIdParams{LocaleCode: localeCode, Id: id})
	if err != nil {
		return nil, err
	}

	result := &profiles.Profile{
		Id:                row.Profile.Id,
		Slug:              row.Profile.Slug,
		Kind:              row.Profile.Kind,
		CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
		ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
		Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
		Title:             row.ProfileTx.Title,
		Description:       row.ProfileTx.Description,
		Properties:        vars.ToRawMessage(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetProfileBySlug(ctx context.Context, localeCode string, slug string) (*profiles.Profile, error) {
	row, err := r.queries.GetProfileBySlug(ctx, GetProfileBySlugParams{LocaleCode: localeCode, Slug: slug})
	if err != nil {
		return nil, err
	}

	result := &profiles.Profile{
		Id:                row.Profile.Id,
		Slug:              row.Profile.Slug,
		Kind:              row.Profile.Kind,
		CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
		ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
		Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
		Title:             row.ProfileTx.Title,
		Description:       row.ProfileTx.Description,
		Properties:        vars.ToRawMessage(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetProfileByCustomDomain(ctx context.Context, localeCode string, domain string) (*profiles.Profile, error) { //nolint:lll
	row, err := r.queries.GetProfileByCustomDomain(ctx, GetProfileByCustomDomainParams{LocaleCode: localeCode, Domain: sql.NullString{String: domain, Valid: true}}) //nolint:lll
	if err != nil {
		return nil, err
	}

	result := &profiles.Profile{
		Id:                row.Profile.Id,
		Slug:              row.Profile.Slug,
		Kind:              row.Profile.Kind,
		CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
		ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
		Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
		Title:             row.ProfileTx.Title,
		Description:       row.ProfileTx.Description,
		Properties:        vars.ToRawMessage(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListProfiles(ctx context.Context, localeCode string) ([]*profiles.Profile, error) {
	rows, err := r.queries.ListProfiles(ctx, ListProfilesParams{LocaleCode: localeCode})
	if err != nil {
		return nil, err
	}

	result := make([]*profiles.Profile, len(rows))
	for i, row := range rows {
		result[i] = &profiles.Profile{
			Id:                row.Profile.Id,
			Slug:              row.Profile.Slug,
			Kind:              row.Profile.Kind,
			CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
			ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
			Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
			Title:             row.ProfileTx.Title,
			Description:       row.ProfileTx.Description,
			Properties:        vars.ToRawMessage(row.Profile.Properties),
			CreatedAt:         row.Profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
		}
	}

	return result, nil
}

func (r *Repository) ListProfilesWithCursor(ctx context.Context, localeCode string, cursor *cursors.Cursor) (cursors.Cursored[[]*profiles.Profile], error) { //nolint:lll
	var wrappedResponse cursors.Cursored[[]*profiles.Profile]

	rows, err := r.queries.ListProfiles(ctx, ListProfilesParams{LocaleCode: localeCode})
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*profiles.Profile, len(rows))
	for i, row := range rows {
		result[i] = &profiles.Profile{
			Id:                row.Profile.Id,
			Slug:              row.Profile.Slug,
			Kind:              row.Profile.Kind,
			CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
			ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
			Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
			Title:             row.ProfileTx.Title,
			Description:       row.ProfileTx.Description,
			Properties:        vars.ToRawMessage(row.Profile.Properties),
			CreatedAt:         row.Profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
		}
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].Id
	}

	return wrappedResponse, nil
}

func (r *Repository) GetProfilePagesByProfileId(ctx context.Context, localeCode string, profileId string) ([]*profiles.ProfilePageBrief, error) { //nolint:lll
	rows, err := r.queries.GetProfilePagesByProfileId(ctx, GetProfilePagesByProfileIdParams{LocaleCode: localeCode, ProfileId: profileId}) //nolint:lll
	if err != nil {
		return nil, err
	}

	profilePages := make([]*profiles.ProfilePageBrief, len(rows))
	for i, row := range rows {
		profilePages[i] = &profiles.ProfilePageBrief{
			Id:              row.Id,
			Slug:            row.Slug,
			CoverPictureUri: vars.ToStringPtr(row.CoverPictureUri),
			Title:           row.Title,
			Summary:         row.Summary,
		}
	}

	return profilePages, nil
}
