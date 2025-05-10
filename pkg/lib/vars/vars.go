package vars

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/sqlc-dev/pqtype"
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

func ToRawMessage(m pqtype.NullRawMessage) *json.RawMessage {
	if m.Valid {
		return &m.RawMessage
	}

	return nil
}
