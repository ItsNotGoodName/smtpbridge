package models

import (
	"database/sql"
	"io"
	"strconv"
	"strings"
	"time"
)

// NOTE: `sql:"primary_key"` tag is used by jet to map database rows to structs

type Config struct {
	RetentionPolicy ConfigRetentionPolicy
	AuthSMTP        Auth
	AuthHTTP        Auth
}

type ConfigRetentionPolicy struct {
	MinAge         time.Duration
	EnvelopeCount  *int
	EnvelopeAge    *time.Duration
	AttachmentSize *int64
	TraceAge       *time.Duration
}

func (p ConfigRetentionPolicy) EnvelopeAgeTime() time.Time {
	if p.EnvelopeAge == nil {
		return time.Time{}
	}

	return time.Now().Add(-*p.EnvelopeAge)
}

func (p ConfigRetentionPolicy) MinAgeTime() time.Time {
	return time.Now().Add(-p.MinAge)
}

type Auth struct {
	Anonymous bool
	Username  string
	Password  string
}

type DataAttachment struct {
	io.Reader
	Attachment Attachment
}

const EndpointTitleTemplate = "{{ .Message.Subject }}"
const EndpointBodyTemplate = "{{ .Message.Text }}"

type Endpoint struct {
	ID                int64 `sql:"primary_key"`
	Internal          bool
	InternalID        sql.NullString
	Name              string
	AttachmentDisable bool
	TextDisable       bool
	TitleTemplate     string
	BodyTemplate      string
	Kind              string
	Config            EndpointConfig
}

type EndpointConfig map[string]string

func (c EndpointConfig) Str(key string) string {
	return string(c[key])
}

func (c EndpointConfig) StrSlice(key string) []string {
	data, _ := c[key]
	return strings.Split(string(data), ",")
}

type Rule struct {
	ID         int64 `sql:"primary_key"`
	Internal   bool
	InternalID sql.NullString
	Name       string
	Expression string
	Enable     bool
}

type RuleEndpoints struct {
	Rule      Rule
	Endpoints []Endpoint
}

type Storage struct {
	EnvelopeCount   int
	AttachmentCount int
	AttachmentSize  int64
	DatabaseSize    int64
}

type Trace struct {
	ID        int64 `sql:"primary_key"`
	Seq       int
	RequestID string
	Source    string
	Action    string
	Level     TraceLevel
	Data      TraceData
	CreatedAt Time
}

type TraceLevel string

type TraceData []TraceDataKV

type TraceDataKV struct {
	Key   string
	Value string
}

func (t TraceDataKV) ValueInt64() int64 {
	id, _ := strconv.ParseInt(t.Value, 10, 64)
	return id
}

type User struct {
	ID       int64 `sql:"primary_key"`
	Username string
}

type EndpointSchemaItem struct {
	Name   string
	Kind   string
	Fields []EndpointSchemaField
}

type EndpointSchemaField struct {
	Name        string
	Description string
	Example     string
	Key         string
	Multiline   bool
	Optional    bool
}

type EndpointSchema []EndpointSchemaItem

func (e EndpointSchema) Get(kind string) EndpointSchemaItem {
	for i := range e {
		if e[i].Kind == kind {
			return e[i]
		}
	}

	return EndpointSchemaItem{}
}

func (e EndpointSchema) Filter(kind string, config EndpointConfig) EndpointConfig {
	for _, end := range e {
		if end.Kind == kind {
			return end.Filter(config)
		}
	}

	return make(EndpointConfig)
}

func (e EndpointSchemaItem) Filter(config EndpointConfig) EndpointConfig {
	validConfig := make(EndpointConfig)

	for _, field := range e.Fields {
		validConfig[field.Key] = config[field.Key]
	}

	return validConfig
}
