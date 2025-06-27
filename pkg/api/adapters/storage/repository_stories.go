//nolint:dupl
package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/eser/aya.is-services/pkg/api/business/profiles"
	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/lib/cursors"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetStoryIDBySlug(ctx context.Context, slug string) (string, error) {
	var result string

	err := r.cache.Execute(
		ctx,
		"story_id_by_slug:"+slug,
		&result,
		func(ctx context.Context) (any, error) {
			row, err := r.queries.GetStoryIDBySlug(ctx, GetStoryIDBySlugParams{Slug: slug})
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, nil //nolint:nilnil
				}

				return nil, err
			}

			return row, nil
		},
	)

	return result, err //nolint:wrapcheck
}

func (r *Repository) GetStoryByID(
	ctx context.Context,
	localeCode string,
	id string,
	authorProfileID *string,
) (*stories.StoryWithChildren, error) {
	getStoryByIDParams := GetStoryByIDParams{
		LocaleCode: localeCode,
		ID:         id,
		FilterPublicationProfileID: sql.NullString{
			String: "",
			Valid:  false,
		},
		FilterAuthorProfileID: sql.NullString{
			String: "",
			Valid:  false,
		},
	}
	if authorProfileID != nil {
		getStoryByIDParams.FilterAuthorProfileID = sql.NullString{
			String: *authorProfileID,
			Valid:  true,
		}
	}

	row, err := r.queries.GetStoryByID(ctx, getStoryByIDParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil //nolint:nilnil
		}

		return nil, err
	}

	result := &stories.StoryWithChildren{
		Story: &stories.Story{
			ID:              row.Story.ID,
			AuthorProfileID: vars.ToStringPtr(row.Story.AuthorProfileID),
			Slug:            row.Story.Slug,
			Kind:            row.Story.Kind,
			Status:          row.Story.Status,
			IsFeatured:      row.Story.IsFeatured,
			StoryPictureURI: vars.ToStringPtr(row.Story.StoryPictureURI),
			Title:           row.StoryTx.Title,
			Summary:         row.StoryTx.Summary,
			Content:         row.StoryTx.Content,
			Properties:      vars.ToObject(row.Story.Properties),
			CreatedAt:       row.Story.CreatedAt,
			UpdatedAt:       vars.ToTimePtr(row.Story.UpdatedAt),
			DeletedAt:       vars.ToTimePtr(row.Story.DeletedAt),
		},
		AuthorProfile: &profiles.Profile{
			ID:                row.Profile.ID,
			Slug:              row.Profile.Slug,
			Kind:              row.Profile.Kind,
			CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
			ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
			Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
			Title:             row.ProfileTx.Title,
			Description:       row.ProfileTx.Description,
			Properties:        vars.ToObject(row.Profile.Properties),
			CreatedAt:         row.Profile.CreatedAt,
			UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
			DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
		},
	}

	return result, nil
}

func (r *Repository) ListStoriesOfPublication( //nolint:funlen
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*stories.StoryWithChildren], error) {
	var wrappedResponse cursors.Cursored[[]*stories.StoryWithChildren]

	rows, err := r.queries.ListStoriesOfPublication(
		ctx,
		ListStoriesOfPublicationParams{
			LocaleCode: localeCode,
			FilterKind: vars.MapValueToNullString(cursor.Filters, "kind"),
			FilterAuthorProfileID: vars.MapValueToNullString(
				cursor.Filters,
				"author_profile_id",
			),
			FilterPublicationProfileID: vars.MapValueToNullString(
				cursor.Filters,
				"publication_profile_id",
			),
		},
	)
	if err != nil {
		return wrappedResponse, err
	}

	result := make([]*stories.StoryWithChildren, len(rows))
	for i, row := range rows {
		storyWithChildren := &stories.StoryWithChildren{
			Story: &stories.Story{
				ID:              row.Story.ID,
				AuthorProfileID: vars.ToStringPtr(row.Story.AuthorProfileID),
				Slug:            row.Story.Slug,
				Kind:            row.Story.Kind,
				Status:          row.Story.Status,
				IsFeatured:      row.Story.IsFeatured,
				StoryPictureURI: vars.ToStringPtr(row.Story.StoryPictureURI),
				Title:           row.StoryTx.Title,
				Summary:         row.StoryTx.Summary,
				Content:         row.StoryTx.Content,
				Properties:      vars.ToObject(row.Story.Properties),
				CreatedAt:       row.Story.CreatedAt,
				UpdatedAt:       vars.ToTimePtr(row.Story.UpdatedAt),
				DeletedAt:       vars.ToTimePtr(row.Story.DeletedAt),
			},
			AuthorProfile: &profiles.Profile{
				ID:                row.Profile.ID,
				Slug:              row.Profile.Slug,
				Kind:              row.Profile.Kind,
				CustomDomain:      vars.ToStringPtr(row.Profile.CustomDomain),
				ProfilePictureURI: vars.ToStringPtr(row.Profile.ProfilePictureURI),
				Pronouns:          vars.ToStringPtr(row.Profile.Pronouns),
				Title:             row.ProfileTx.Title,
				Description:       row.ProfileTx.Description,
				Properties:        vars.ToObject(row.Profile.Properties),
				CreatedAt:         row.Profile.CreatedAt,
				UpdatedAt:         vars.ToTimePtr(row.Profile.UpdatedAt),
				DeletedAt:         vars.ToTimePtr(row.Profile.DeletedAt),
			},
			Publications: nil,
		}

		var publicationProfiles []struct {
			Profile struct {
				CreatedAt         time.Time        `db:"created_at"          json:"created_at"`
				CustomDomain      *string          `db:"custom_domain"       json:"custom_domain"`
				ProfilePictureURI *string          `db:"profile_picture_uri" json:"profile_picture_uri"`
				Pronouns          *string          `db:"pronouns"            json:"pronouns"`
				Properties        *json.RawMessage `db:"properties"          json:"properties"`
				UpdatedAt         *time.Time       `db:"updated_at"          json:"updated_at"`
				DeletedAt         *time.Time       `db:"deleted_at"          json:"deleted_at"`
				ID                string           `db:"id"                  json:"id"`
				Slug              string           `db:"slug"                json:"slug"`
				Kind              string           `db:"kind"                json:"kind"`
			} `json:"profile"`
			ProfileTx struct {
				Properties  *json.RawMessage `db:"properties"  json:"properties"`
				ProfileID   string           `db:"profile_id"  json:"profile_id"`
				LocaleCode  string           `db:"locale_code" json:"locale_code"`
				Title       string           `db:"title"       json:"title"`
				Description string           `db:"description" json:"description"`
			} `json:"profile_tx"`
		}

		err := json.Unmarshal(row.Publications, &publicationProfiles)
		if err != nil {
			r.logger.Error("failed to unmarshal publications", "error", err)

			continue
		}

		storyWithChildren.Publications = make([]*profiles.Profile, len(publicationProfiles))
		for j, publicationProfile := range publicationProfiles {
			storyWithChildren.Publications[j] = &profiles.Profile{
				ID:                publicationProfile.Profile.ID,
				Slug:              publicationProfile.Profile.Slug,
				Kind:              publicationProfile.Profile.Kind,
				CustomDomain:      publicationProfile.Profile.CustomDomain,
				ProfilePictureURI: publicationProfile.Profile.ProfilePictureURI,
				Pronouns:          publicationProfile.Profile.Pronouns,
				Title:             publicationProfile.ProfileTx.Title,
				Description:       publicationProfile.ProfileTx.Description,
				Properties:        publicationProfile.Profile.Properties,
				CreatedAt:         publicationProfile.Profile.CreatedAt,
				UpdatedAt:         publicationProfile.Profile.UpdatedAt,
				DeletedAt:         publicationProfile.Profile.DeletedAt,
			}
		}

		result[i] = storyWithChildren
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}
