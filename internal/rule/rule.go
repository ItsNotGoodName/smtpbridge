package rule

import (
	"bytes"
	"database/sql"
	"fmt"
	"strconv"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/internal/models"
)

func validate(r models.Rule) error {
	if r.Internal && !r.InternalID.Valid {
		return fmt.Errorf("internal id is empty")
	}

	if r.Name == "" {
		return models.FieldError{Field: models.FieldName, Err: fmt.Errorf("cannot be empty")}
	}

	t, err := TemplateBuild(r.Expression)
	if err != nil {
		return models.FieldError{Field: models.FieldExpression, Err: err}
	}

	_, err = TemplateRun(t, models.Envelope{})
	if err != nil {
		return models.FieldError{Field: models.FieldExpression, Err: err}
	}

	return nil
}

func new(r models.DTORuleCreate) models.Rule {
	return models.Rule{
		Internal:   false,
		InternalID: sql.NullString{},
		Name:       r.Name,
		Expression: r.Expression,
		Enable:     true,
	}
}

func New(r models.DTORuleCreate) (models.Rule, error) {
	rule := new(r)
	return rule, validate(rule)
}

func NewInternal(r models.DTORuleCreate, internalID string) (models.Rule, error) {
	if r.Name == "" {
		r.Name = internalID
	}

	rule := new(r)
	rule.Internal = true
	rule.InternalID = sql.NullString{
		String: internalID,
		Valid:  true,
	}

	return rule, validate(rule)
}

func Update(r models.Rule, req models.DTORuleUpdate) (models.Rule, error) {
	if r.Internal {
		if req.Name == nil && req.Expression == nil && req.Enable != nil {
			r.Enable = *req.Enable
			return r, validate(r)
		}

		return models.Rule{}, models.ErrInternalResource
	}

	if req.Name != nil {
		r.Name = *req.Name
	}

	if req.Expression != nil {
		r.Expression = *req.Expression
	}

	if req.Enable != nil {
		r.Enable = *req.Enable
	}

	return r, validate(r)
}

func Delete(r models.Rule) error {
	if r.Internal {
		return models.ErrInternalResource
	}
	return nil
}

func TemplateBuild(expression string) (*template.Template, error) {
	t := expression
	if t == "" {
		t = "\"true\""
	}

	return template.New("").Parse("{{" + t + "}}")
}

func TemplateRun(tmpl *template.Template, env models.Envelope) (bool, error) {
	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, env)
	if err != nil {
		return false, err
	}

	return strconv.ParseBool(buffer.String())
}

type Rule struct {
	ID       int64
	Template *template.Template
}

func Build(r models.Rule) (Rule, error) {
	tmpl, err := TemplateBuild(r.Expression)
	if err != nil {
		return Rule{}, err
	}

	return Rule{
		ID:       r.ID,
		Template: tmpl,
	}, nil
}

func (r Rule) Match(env models.Envelope) (bool, error) {
	return TemplateRun(r.Template, env)
}
