package rules

import (
	"bytes"
	"fmt"
	"strconv"
	"text/template"

	"github.com/ItsNotGoodName/smtpbridge/internal/envelope"
)

type ParsedRule struct {
	ID       int64
	Template *template.Template
}

type CreateRule struct {
	Internal   bool
	InternalID string
	Name       string
	Template   string
	Enable     bool
}

type Rule struct {
	ID         int64 `bun:"id,pk,autoincrement"`
	Internal   bool
	InternalID string
	Name       string
	Template   string
	Enable     bool
}

func New(req CreateRule) (Rule, error) {
	if req.Internal && req.InternalID == "" {
		return Rule{}, fmt.Errorf("internal id is empty")
	}

	if req.Name == "" {
		return Rule{}, fmt.Errorf("name is empty")
	}

	rrule := Rule{
		Internal:   req.Internal,
		InternalID: req.InternalID,
		Name:       req.Name,
		Template:   req.Template,
		Enable:     true,
	}

	_, err := rrule.Parse()
	return rrule, err
}

func (r Rule) Parse() (ParsedRule, error) {
	t := r.Template
	if t == "" {
		t = "\"true\""
	}
	tmpl, err := template.New("").Parse("{{" + t + "}}")
	if err != nil {
		return ParsedRule{}, err
	}

	return ParsedRule{
		ID:       r.ID,
		Template: tmpl,
	}, nil
}

func (pr ParsedRule) Match(env envelope.Envelope) (bool, error) {
	var buffer bytes.Buffer
	if err := pr.Template.Execute(&buffer, env); err != nil {
		return false, err
	}

	truthy, err := strconv.ParseBool(buffer.String())
	if err != nil {
		return false, err
	}

	return truthy, nil
}

type Aggregate struct {
	Rule      Rule
	Endpoints []Endpoint
}

type Endpoint struct {
	ID     int64
	Name   string
	Enable bool
}
