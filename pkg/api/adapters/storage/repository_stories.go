//nolint:dupl
package storage

import (
	"context"

	"github.com/eser/aya.is-services/pkg/api/business/stories"
	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (r *Repository) GetStoryById(ctx context.Context, localeCode string, id string) (*stories.Story, error) {
	row, err := r.queries.GetStoryById(ctx, GetStoryByIdParams{LocaleCode: localeCode, Id: id})
	if err != nil {
		return nil, err
	}

	result := &stories.Story{
		Id:              row.Story.Id,
		AuthorProfileId: vars.ToStringPtr(row.Story.AuthorProfileId),
		Slug:            row.Story.Slug,
		Kind:            row.Story.Kind,
		Status:          row.Story.Status,
		IsFeatured:      row.Story.IsFeatured,
		StoryPictureUri: vars.ToStringPtr(row.Story.StoryPictureUri),
		Title:           row.StoryTx.Title,
		Summary:         row.StoryTx.Summary,
		Content:         row.StoryTx.Content,
		Properties:      vars.ToRawMessage(row.Story.Properties),
		PublishedAt:     vars.ToTimePtr(row.Story.PublishedAt),
		CreatedAt:       row.Story.CreatedAt,
		UpdatedAt:       vars.ToTimePtr(row.Story.UpdatedAt),
		DeletedAt:       vars.ToTimePtr(row.Story.DeletedAt),
	}

	return result, nil
}

func (r *Repository) GetStoryBySlug(ctx context.Context, localeCode string, slug string) (*stories.Story, error) {
	row, err := r.queries.GetStoryBySlug(ctx, GetStoryBySlugParams{LocaleCode: localeCode, Slug: slug})
	if err != nil {
		return nil, err
	}

	result := &stories.Story{
		Id:              row.Story.Id,
		AuthorProfileId: vars.ToStringPtr(row.Story.AuthorProfileId),
		Slug:            row.Story.Slug,
		Kind:            row.Story.Kind,
		Status:          row.Story.Status,
		IsFeatured:      row.Story.IsFeatured,
		StoryPictureUri: vars.ToStringPtr(row.Story.StoryPictureUri),
		Title:           row.StoryTx.Title,
		Summary:         row.StoryTx.Summary,
		Content:         row.StoryTx.Content,
		Properties:      vars.ToRawMessage(row.Story.Properties),
		PublishedAt:     vars.ToTimePtr(row.Story.PublishedAt),
		CreatedAt:       row.Story.CreatedAt,
		UpdatedAt:       vars.ToTimePtr(row.Story.UpdatedAt),
		DeletedAt:       vars.ToTimePtr(row.Story.DeletedAt),
	}

	return result, nil
}

func (r *Repository) ListStories(ctx context.Context, localeCode string) ([]*stories.Story, error) {
	rows, err := r.queries.ListStories(ctx, ListStoriesParams{LocaleCode: localeCode})
	if err != nil {
		return nil, err
	}

	result := make([]*stories.Story, len(rows))
	for i, row := range rows {
		result[i] = &stories.Story{
			Id:              row.Story.Id,
			AuthorProfileId: vars.ToStringPtr(row.Story.AuthorProfileId),
			Slug:            row.Story.Slug,
			Kind:            row.Story.Kind,
			Status:          row.Story.Status,
			IsFeatured:      row.Story.IsFeatured,
			StoryPictureUri: vars.ToStringPtr(row.Story.StoryPictureUri),
			Title:           row.StoryTx.Title,
			Summary:         row.StoryTx.Summary,
			Content:         row.StoryTx.Content,
			Properties:      vars.ToRawMessage(row.Story.Properties),
			PublishedAt:     vars.ToTimePtr(row.Story.PublishedAt),
			CreatedAt:       row.Story.CreatedAt,
			UpdatedAt:       vars.ToTimePtr(row.Story.UpdatedAt),
			DeletedAt:       vars.ToTimePtr(row.Story.DeletedAt),
		}
	}

	return result, nil
}
