package models

type Field string

type FieldError struct {
	Field Field
	Err   error
}

func (e FieldError) Error() string {
	return "FieldError" + ": " + string(e.Field) + ": " + e.Err.Error()
}

func (e FieldError) Unwrap() error {
	return e.Err
}

const (
	FieldName          Field = "name"
	FieldExpression    Field = "expression"
	FieldTitleTemplate Field = "title template"
	FieldBodyTemplate  Field = "body template"
	FieldKind          Field = "kind"
	FieldConfig        Field = "config"
)
