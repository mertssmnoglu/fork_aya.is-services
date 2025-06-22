package profiles

import (
	"time"

	"github.com/eser/aya.is-services/pkg/ajan/lib"
)

type RecordID string

type RecordIDGenerator func() RecordID

func DefaultIDGenerator() RecordID {
	return RecordID(lib.IDsGenerateUnique())
}

type Profile struct {
	CreatedAt         time.Time  `json:"created_at"`
	Properties        any        `json:"properties"`
	CustomDomain      *string    `json:"custom_domain"`
	ProfilePictureURI *string    `json:"profile_picture_uri"`
	Pronouns          *string    `json:"pronouns"`
	UpdatedAt         *time.Time `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at"`
	ID                string     `json:"id"`
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
	CoverPictureURI *string    `json:"cover_picture_uri"`
	PublishedAt     *time.Time `json:"published_at"`
	ID              string     `json:"id"`
	Slug            string     `json:"slug"`
	Title           string     `json:"title"`
	Summary         string     `json:"summary"`
	Content         string     `json:"content"`
}

type ProfilePageBrief struct {
	ID              string  `json:"id"`
	Slug            string  `json:"slug"`
	CoverPictureURI *string `json:"cover_picture_uri"`
	Title           string  `json:"title"`
	Summary         string  `json:"summary"`
}

type ProfileLinkBrief struct {
	ID         string `json:"id"`
	Kind       string `json:"kind"`
	PublicID   string `json:"public_id"`
	URI        string `json:"uri"`
	Title      string `json:"title"`
	IsVerified bool   `json:"is_verified"`
}

type ProfileMembership struct {
	Properties    any        `json:"properties"`
	Profile       *Profile   `json:"profile"`
	MemberProfile *Profile   `json:"member_profile"`
	StartedAt     *time.Time `json:"started_at"`
	FinishedAt    *time.Time `json:"finished_at"`
	ID            string     `json:"id"`
	Kind          string     `json:"kind"`
}

type ExternalPost struct {
	CreatedAt *time.Time `json:"created_at"` //nolint:tagliatelle
	ID        string     `json:"id"`
	Content   string     `json:"content"`
	Permalink string     `json:"permalink"`
}
