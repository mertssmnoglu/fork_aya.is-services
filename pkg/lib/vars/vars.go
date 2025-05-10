package vars

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/sqlc-dev/pqtype"
)

var (
	ErrMustBePointer  = errors.New("must be a pointer")
	ErrCannotAssign   = errors.New("cannot assign")
	ErrCannotSetValue = errors.New("cannot set a value")
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

func ToRawMessage(m pqtype.NullRawMessage) []byte {
	if m.Valid {
		return m.RawMessage
	}

	return nil
}

func SetValue(dest any, src any) error {
	dv := reflect.ValueOf(dest)
	if dv.Kind() != reflect.Pointer {
		return fmt.Errorf("%w: %s", ErrMustBePointer, dv.Kind())
	}

	ev := dv.Elem() //nolint:varnamelen

	sv := reflect.ValueOf(src) //nolint:varnamelen
	if !sv.Type().AssignableTo(ev.Type()) {
		return fmt.Errorf("%w: %s to %s", ErrCannotAssign, sv.Type(), ev.Type())
	}

	if !ev.CanSet() {
		return fmt.Errorf("%w: %s", ErrCannotSetValue, ev.Type())
	}

	ev.Set(sv)

	return nil
}
