// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: query.sql

package queries

import (
	"context"
	"time"
)

const deleteEnvelopeOlderThan = `-- name: DeleteEnvelopeOlderThan :execrows
DELETE FROM messages WHERE created_at < ?
`

func (q *Queries) DeleteEnvelopeOlderThan(ctx context.Context, createdAt time.Time) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEnvelopeOlderThan, createdAt)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const deleteEnvelopeUntilCount = `-- name: DeleteEnvelopeUntilCount :execrows
DELETE FROM messages
WHERE id NOT IN (
  SELECT id
  FROM messages
  ORDER BY id DESC
  LIMIT ?
) AND messages.created_at < ?
`

type DeleteEnvelopeUntilCountParams struct {
	Limit     int64
	CreatedAt time.Time
}

func (q *Queries) DeleteEnvelopeUntilCount(ctx context.Context, arg DeleteEnvelopeUntilCountParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, deleteEnvelopeUntilCount, arg.Limit, arg.CreatedAt)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const getEnvelopeMessageHTML = `-- name: GetEnvelopeMessageHTML :one
SELECT html FROM messages WHERE id = ?1
`

func (q *Queries) GetEnvelopeMessageHTML(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getEnvelopeMessageHTML, id)
	var html string
	err := row.Scan(&html)
	return html, err
}

const isRuleInternal = `-- name: IsRuleInternal :one
;

SELECT internal FROM rules WHERE id = ?1
`

func (q *Queries) IsRuleInternal(ctx context.Context, id int64) (bool, error) {
	row := q.db.QueryRowContext(ctx, isRuleInternal, id)
	var internal bool
	err := row.Scan(&internal)
	return internal, err
}

const upsertInternalEndpoint = `-- name: UpsertInternalEndpoint :exec
INSERT INTO endpoints (
  internal,
  internal_id,
  name,
  attachment_disable,
  text_disable,
  body_template,
  kind,
  config,
  updated_at
) VALUES (
  true,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?,
  ?
) ON CONFLICT (internal_id) DO UPDATE SET
  internal_id=EXCLUDED.internal_id,
  name=EXCLUDED.name,
  attachment_disable=EXCLUDED.attachment_disable,
  text_disable=EXCLUDED.text_disable,
  body_template=EXCLUDED.body_template,
  kind=EXCLUDED.kind,
  config=EXCLUDED.config,
  updated_at=EXCLUDED.updated_at WHERE internal = true
`

type UpsertInternalEndpointParams struct {
	InternalID        string
	Name              string
	AttachmentDisable bool
	TextDisable       bool
	BodyTemplate      string
	Kind              string
	Config            interface{}
	UpdatedAt         time.Time
}

// autocmd BufWritePost query.sql !sqlc generate
func (q *Queries) UpsertInternalEndpoint(ctx context.Context, arg UpsertInternalEndpointParams) error {
	_, err := q.db.ExecContext(ctx, upsertInternalEndpoint,
		arg.InternalID,
		arg.Name,
		arg.AttachmentDisable,
		arg.TextDisable,
		arg.BodyTemplate,
		arg.Kind,
		arg.Config,
		arg.UpdatedAt,
	)
	return err
}

const upsertInternalRule = `-- name: UpsertInternalRule :exec
INSERT INTO rules (
  internal,
  internal_id,
  name,
  expression,
  updated_at,
  enable
) VALUES (
  true,
  ?,
  ?,
  ?,
  ?,
  ?
) ON CONFLICT (internal_id) DO UPDATE SET
  internal_id=EXCLUDED.internal_id,
  name=EXCLUDED.name,
  expression=EXCLUDED.expression,
  updated_at=EXCLUDED.updated_at
WHERE internal = true
`

type UpsertInternalRuleParams struct {
	InternalID string
	Name       string
	Expression string
	UpdatedAt  time.Time
	Enable     bool
}

func (q *Queries) UpsertInternalRule(ctx context.Context, arg UpsertInternalRuleParams) error {
	_, err := q.db.ExecContext(ctx, upsertInternalRule,
		arg.InternalID,
		arg.Name,
		arg.Expression,
		arg.UpdatedAt,
		arg.Enable,
	)
	return err
}

const upsertInternalRuleToEndpoint = `-- name: UpsertInternalRuleToEndpoint :execrows
INSERT INTO rules_to_endpoints (
  internal,
  rule_id,
  endpoint_id,
  updated_at
) SELECT "1" AS internal, rules.id AS rule_id, endpoints.id AS endpoint_id, ?
FROM rules, endpoints
WHERE rules.internal_id=? AND endpoints.internal_id IN (?) 
ON CONFLICT (rule_id, endpoint_id) DO UPDATE SET updated_at=EXCLUDED.updated_at, internal=EXCLUDED.internal
`

type UpsertInternalRuleToEndpointParams struct {
	UpdatedAt    time.Time
	InternalID   string
	InternalID_2 string
}

func (q *Queries) UpsertInternalRuleToEndpoint(ctx context.Context, arg UpsertInternalRuleToEndpointParams) (int64, error) {
	result, err := q.db.ExecContext(ctx, upsertInternalRuleToEndpoint, arg.UpdatedAt, arg.InternalID, arg.InternalID_2)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
