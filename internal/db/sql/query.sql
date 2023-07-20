-- autocmd BufWritePost query.sql !sqlc generate
-- name: UpsertInternalEndpoint :exec
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
  updated_at=EXCLUDED.updated_at WHERE internal = true;

-- name: UpsertInternalRule :exec
INSERT INTO rules (
  internal,
  internal_id,
  name,
  template,
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
  template=EXCLUDED.template,
  updated_at=EXCLUDED.updated_at
WHERE internal = true;

-- name: UpsertInternalRuleToEndpoint :execrows
INSERT INTO rules_to_endpoints (
  internal,
  rule_id,
  endpoint_id,
  updated_at
) SELECT "1" AS internal, rules.id AS rule_id, endpoints.id AS endpoint_id, ?
FROM rules, endpoints
WHERE rules.internal_id=? AND endpoints.internal_id IN (?) 
ON CONFLICT (rule_id, endpoint_id) DO UPDATE SET updated_at=EXCLUDED.updated_at, internal=EXCLUDED.internal;

-- name: GetEnvelopeMessageHTML :one
SELECT html FROM messages WHERE id = ?1;

-- name: DeleteEnvelopeUntilCount :execrows
DELETE FROM messages
WHERE id NOT IN (
  SELECT id
  FROM messages
  ORDER BY id DESC
  LIMIT ?
);

-- name: DeleteEnvelopeOlderThan :execrows
DELETE FROM messages WHERE created_at < ? ;

-- name: IsRuleInternal :one
SELECT internal FROM rules WHERE id = ?1;
