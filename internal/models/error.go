package models

import "fmt"

type ErrorField struct {
	Field   string
	Message string
}

var (
	ErrNotFound    = fmt.Errorf("not found")
	ErrAuthInvalid = fmt.Errorf("auth invalid")
)
