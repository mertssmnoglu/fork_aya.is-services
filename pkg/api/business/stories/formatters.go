package stories

import (
	"encoding/json"
	"time"

	"github.com/eser/aya.is-services/pkg/lib/vars"
)

func (s *Story) MarshalJSON() ([]byte, error) {
	formatted := &struct {
		Id              string     `json:"id"`
		AuthorProfileId *string    `json:"author_profile_id"`
		Slug            string     `json:"slug"`
		Kind            string     `json:"kind"`
		Status          string     `json:"status"`
		IsFeatured      bool       `json:"is_featured"`
		StoryPictureUri *string    `json:"story_picture_uri"`
		Title           string     `json:"title"`
		Summary         string     `json:"summary"`
		Content         string     `json:"content"`
		PublishedAt     *time.Time `json:"published_at"`
		CreatedAt       time.Time  `json:"created_at"`
		UpdatedAt       *time.Time `json:"updated_at"`
		DeletedAt       *time.Time `json:"deleted_at"`
	}{
		Id:              s.Id,
		AuthorProfileId: vars.ToStringPtr(s.AuthorProfileId),
		Slug:            s.Slug,
		Kind:            s.Kind,
		Status:          s.Status,
		IsFeatured:      s.IsFeatured,
		StoryPictureUri: vars.ToStringPtr(s.StoryPictureUri),
		Title:           s.Title,
		Summary:         s.Summary,
		Content:         s.Content,
		PublishedAt:     vars.ToTimePtr(s.PublishedAt),
		CreatedAt:       s.CreatedAt,
		UpdatedAt:       vars.ToTimePtr(s.UpdatedAt),
		DeletedAt:       vars.ToTimePtr(s.DeletedAt),
	}

	return json.Marshal(formatted) //nolint:wrapcheck
}
