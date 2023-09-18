package models

import "fmt"

var (
	ErrNotFound         = fmt.Errorf("not found")
	ErrAuthInvalid      = fmt.Errorf("auth invalid")
	ErrInternalResource = fmt.Errorf("internal resource")
)
