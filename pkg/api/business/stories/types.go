package stories

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(ulid.Make().String())
}

type Story struct {
	CreatedAt       time.Time  `json:"created_at"`
	Properties      any        `json:"properties"`
	AuthorProfileId *string    `json:"author_profile_id"`
	StoryPictureUri *string    `json:"story_picture_uri"`
	PublishedAt     *time.Time `json:"published_at"`
	UpdatedAt       *time.Time `json:"updated_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	Id              string     `json:"id"`
	Slug            string     `json:"slug"`
	Kind            string     `json:"kind"`
	Status          string     `json:"status"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
	IsFeatured      bool       `json:"is_featured"`
}
