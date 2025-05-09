package vars

import (
	"database/sql"
	"time"
)

func ToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}

	return nil
}

func ToTimePtr(t sql.NullTime) *time.Time {
	if t.Valid {
		return &t.Time
	}

	return nil
}
