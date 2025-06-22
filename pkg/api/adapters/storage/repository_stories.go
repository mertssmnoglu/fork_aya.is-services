//nolint:dupl
package storage

import (
	"context"
	"database/sql"
	"errors"

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
) (*stories.StoryWithChildren, error) {
	row, err := r.queries.GetStoryByID(ctx, GetStoryByIDParams{LocaleCode: localeCode, ID: id})
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
			PublishedAt:     vars.ToTimePtr(row.Story.PublishedAt),
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

func (r *Repository) ListStories( //nolint:funlen
	ctx context.Context,
	localeCode string,
	cursor *cursors.Cursor,
) (cursors.Cursored[[]*stories.StoryWithChildren], error) {
	var wrappedResponse cursors.Cursored[[]*stories.StoryWithChildren]

	rows, err := r.queries.ListStories(
		ctx,
		ListStoriesParams{
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
		result[i] = &stories.StoryWithChildren{
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
				PublishedAt:     vars.ToTimePtr(row.Story.PublishedAt),
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
	}

	wrappedResponse.Data = result

	if len(result) == cursor.Limit {
		wrappedResponse.CursorPtr = &result[len(result)-1].ID
	}

	return wrappedResponse, nil
}
