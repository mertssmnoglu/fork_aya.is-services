package profiles

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(ulid.Make().String())
}

type Profile struct {
	CreatedAt         time.Time  `json:"created_at"`
	Properties        any        `json:"properties"`
	CustomDomain      *string    `json:"custom_domain"`
	ProfilePictureUri *string    `json:"profile_picture_uri"`
	Pronouns          *string    `json:"pronouns"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	Id                string     `json:"id"`
	Slug              string     `json:"slug"`
	Kind              string     `json:"kind"`
	Title             string     `json:"title"`
	Description       string     `json:"description"`
}

type ProfileWithChildren struct {
	*Profile
	Pages []*ProfilePageBrief `json:"pages"`
	Links []*ProfileLinkBrief `json:"links"`
}

type ProfilePage struct {
	CoverPictureUri *string    `json:"cover_picture_uri"`
	PublishedAt     *time.Time `json:"published_at"`
	Id              string     `json:"id"`
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
}

type ProfilePageBrief struct {
	Id              string  `json:"id"`
	Slug            string  `json:"slug"`
	CoverPictureUri *string `json:"cover_picture_uri"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
}

type ProfileLinkBrief struct {
	Id         string `json:"id"`
	Kind       string `json:"kind"`
	PublicId   string `json:"public_id"`
	Uri        string `json:"uri"`
	Title      string `json:"title"`
	IsVerified bool   `json:"is_verified"`
}

type ProfileMembership struct {
	Properties any      `json:"properties"`
	Profile    *Profile `json:"profile"`
	Id         string   `json:"id"`
	Kind       string   `json:"kind"`
}

type ExternalPost struct {
	CreatedAt *time.Time `json:"created_at"` //nolint:tagliatelle
	Id        string     `json:"id"`
	Content   string     `json:"content"`
	Permalink string     `json:"permalink"`
}
