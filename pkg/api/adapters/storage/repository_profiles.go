package storage

import (
	"context"
	"database/sql"
	"strings"

	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetProfileIdBySlug(ctx context.Context, slug string) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"profile_id_by_slug:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetProfileIdBySlug(ctx, GetProfileIdBySlugParams{Slug: slug})
			if err != nil {
				return nil, err
			}

			return row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetProfileIdByCustomDomain(
	ctx context.Context,
	domain string,
) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"profile_id_by_custom_domain:"+domain,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetProfileIdByCustomDomain(
				ctx,
				GetProfileIdByCustomDomainParams{
					CustomDomain: sql.NullString{String: domain, Valid: true},
				},
			)
			if err != nil {
				return nil, err
			}

			return row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetProfileById(
	ctx context.Context,
	localeCode string,
	id string,
) (*profiles.Profile, error) {
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
		Properties:        vars.ToObject(row.Profile.Properties),
		CreatedAt:         row.Profile.CreatedAt,
		UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
		DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListProfiles(
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.Profile], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.Profile]

	rows, err := r.queries.ListProfiles(
		ctx,
		ListProfilesParams{
			LocaleCode: localeCode,
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
		},
	)
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
			Properties:        vars.ToObject(row.Profile.Properties),
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

func (r *Repository) GetProfilePagesByProfileId(
	ctx context.Context,
	localeCode string,
	profileId string,
) ([]*profiles.ProfilePageBrief, error) {
	rows, err := r.queries.GetProfilePagesByProfileId(
		ctx,
		GetProfilePagesByProfileIdParams{LocaleCode: localeCode, ProfileId: profileId},
	)
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

func (r *Repository) GetProfilePageByProfileIdAndSlug(
	ctx context.Context,
	localeCode string,
	profileId string,
	pageSlug string,
) (*profiles.ProfilePage, error) {
	row, err := r.queries.GetProfilePageByProfileIdAndSlug(
		ctx,
		GetProfilePageByProfileIdAndSlugParams{
			LocaleCode: localeCode,
			ProfileId:  profileId,
			PageSlug:   pageSlug,
		},
	)
	if err != nil {
		return nil, err
	}

	result := &profiles.ProfilePage{
		Id:              row.Id,
		Slug:            row.Slug,
		CoverPictureUri: vars.ToStringPtr(row.CoverPictureUri),
		Title:           row.Title,
		Summary:         row.Summary,
		Content:         row.Content,
		PublishedAt:     vars.ToTimePtr(row.PublishedAt),
	}

	return result, nil
}

func (r *Repository) GetProfileLinksByProfileId(
	ctx context.Context,
	_localeCode string,
	profileId string,
) ([]*profiles.ProfileLinkBrief, error) {
	rows, err := r.queries.GetProfileLinksByProfileId(
		ctx,
		GetProfileLinksByProfileIdParams{ProfileId: profileId},
	)
	if err != nil {
		return nil, err
	}

	profileLinks := make([]*profiles.ProfileLinkBrief, len(rows))
	for i, row := range rows {
		profileLinks[i] = &profiles.ProfileLinkBrief{
			Id:         row.Id,
			Kind:       row.Kind,
			IsVerified: row.IsVerified,
			PublicId:   row.PublicId.String,
			Uri:        row.Uri.String,
			Title:      row.Title,
		}
	}

	return profileLinks, nil
}

func (r *Repository) ListProfileContributions(
	ctx context.Context,
	localeCode string,
	profileId string,
	kinds []string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*profiles.ProfileMembership], error) {
	var wrappedResponse cursors.Cursored[[]*profiles.ProfileMembership]

	rows, err := r.queries.ListProfileMemberships(
		ctx,
		ListProfileMembershipsParams{
			LocaleCode:              localeCode,
			FilterProfileId:         sql.NullString{String: "", Valid: false},
			FilterProfileKind:       sql.NullString{String: strings.Join(kinds, ","), Valid: true},
			FilterMemberProfileId:   sql.NullString{String: profileId, Valid: true},
			FilterMemberProfileKind: sql.NullString{String: "", Valid: false},
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	profileMemberships := make([]*profiles.ProfileMembership, len(rows))
	for i, row := range rows {
		profileMemberships[i] = &profiles.ProfileMembership{
			Id:         row.ProfileMembership.Id,
			Kind:       row.ProfileMembership.Kind,
			StartedAt:  vars.ToTimePtr(row.ProfileMembership.StartedAt),
			FinishedAt: vars.ToTimePtr(row.ProfileMembership.FinishedAt),
			Properties: vars.ToObject(row.ProfileMembership.Properties),
			Profile: &profiles.Profile{
				Id:                row.Profile.Id,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
				ProfilePictureUri: vars.ToStringPtr(row.Profile.ProfilePictureUri),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			MemberProfile: &profiles.Profile{
				Id:                row.Profile_2.Id,
				Slug:              row.Profile_2.Slug,
				Kind:              row.Profile_2.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile_2.CustomDomain),
				ProfilePictureUri: vars.ToStringPtr(row.Profile_2.ProfilePictureUri),
				Pronouns:          vars.ToStringPtr(row.Profile_2.Pronouns),
				Title:             row.ProfileTx_2.Title,
				Description:       row.ProfileTx_2.Description,
				Properties:        vars.ToObject(row.Profile_2.Properties),
				CreatedAt:         row.Profile_2.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile_2.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile_2.DeletedAt),
			},
		}
	}

	wrappedResponse.Data = profileMemberships

	if len(profileMemberships) == cursor.Limit {
		wrappedResponse.CursorPtr = &profileMemberships[len(profileMemberships)-1].Id
	}

	return wrappedResponse, nil
}
