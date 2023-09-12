package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

func (dst *TraceData) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		return json.Unmarshal([]byte(src), &dst)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src TraceData) Value() (driver.Value, error) {
	if src == nil {
		return "[]", nil
	}

	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (dst *EndpointConfig) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		return json.Unmarshal([]byte(src), &dst)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src EndpointConfig) Value() (driver.Value, error) {
	if src == nil {
		return "{}", nil
	}

	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

func (dst *StringSlice) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case string:
		return json.Unmarshal([]byte(src), &dst)
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src StringSlice) Value() (driver.Value, error) {
	if src == nil {
		return "[]", nil
	}
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return string(data), nil
}

type Time time.Time

func NewTime(time time.Time) Time {
	return Time(time.UTC())
}

func (t Time) Time() time.Time {
	return time.Time(t)
}

func (dst *Time) Scan(src any) error {
	if src == nil {
		return fmt.Errorf("cannot scan nil")
	}

	switch src := src.(type) {
	case time.Time:
		*dst = Time(src)
		return nil
	}

	return fmt.Errorf("cannot scan %T", src)
}

func (src Time) Value() (driver.Value, error) {
	return time.Time(src).UTC().Format(time.RFC3339), nil
}
