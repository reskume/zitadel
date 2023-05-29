package database

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/jackc/pgtype"
)

type TextArray[t ~string] []t

// Scan implements the [database/sql.Scanner] interface.
func (s *TextArray[t]) Scan(src any) error {
	array := new(pgtype.TextArray)
	if err := array.Scan(src); err != nil {
		return err
	}
	return array.AssignTo(s)
}

// Value implements the [database/sql/driver.Valuer] interface.
func (s TextArray[t]) Value() (driver.Value, error) {
	if len(s) == 0 {
		return nil, nil
	}

	array := pgtype.TextArray{}
	if err := array.Set(s); err != nil {
		return nil, err
	}

	return array.Value()
}

type arrayField interface {
	~int8 | ~uint8 | ~int16 | ~uint16 | ~int32 | ~uint32
}

type Array[F arrayField] []F

// Scan implements the [database/sql.Scanner] interface.
func (a *Array[F]) Scan(src any) error {
	array := new(pgtype.Int8Array)
	if err := array.Scan(src); err != nil {
		return err
	}
	elements := make([]int64, len(array.Elements))
	if err := array.AssignTo(&elements); err != nil {
		return err
	}
	*a = make([]F, len(elements))
	for i, element := range elements {
		(*a)[i] = F(element)
	}
	return nil
}

// Value implements the [database/sql/driver.Valuer] interface.
func (a Array[F]) Value() (driver.Value, error) {
	if len(a) == 0 {
		return nil, nil
	}

	array := pgtype.Int8Array{}
	if err := array.Set(a); err != nil {
		return nil, err
	}

	return array.Value()
}

type Map[V any] map[string]V

// Scan implements the [database/sql.Scanner] interface.
func (m *Map[V]) Scan(src any) error {
	bytea := new(pgtype.Bytea)
	if err := bytea.Scan(src); err != nil {
		return err
	}
	if len(bytea.Bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytea.Bytes, &m)
}

// Value implements the [database/sql/driver.Valuer] interface.
func (m Map[V]) Value() (driver.Value, error) {
	if len(m) == 0 {
		return nil, nil
	}
	return json.Marshal(m)
}
