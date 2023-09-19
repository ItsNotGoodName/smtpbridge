table "messages" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "uuid" {
    null = false
    type = text
    default = ""
  }
  column "from" {
    null = false
    type = text
  }
  column "to" {
    null = false
    type = json
  }
  column "subject" {
    null = false
    type = text
  }
  column "text" {
    null = false
    type = text
  }
  column "html" {
    null = false
    type = text
  }
  column "date" {
    null = false
    type = datetime
  }
  column "created_at" {
    null    = false
    type    = datetime
  }
  primary_key {
    columns = [column.id]
  }
}
table "attachments" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "uuid" {
    null = false
    type = text
    default = ""
  }
  column "message_id" {
    null = true
    type = integer
  }
  column "name" {
    null = false
    type = text
  }
  column "mime" {
    null = false
    type = text
  }
  column "extension" {
    null = false
    type = text
  }
  primary_key {
    columns = [column.id]
  }
  foreign_key "message_id" {
    columns     = [column.message_id]
    ref_columns = [table.messages.column.id]
    on_update   = CASCADE
    on_delete   = SET_NULL
  }
}
table "endpoints" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "internal" {
    null = false
    type = boolean
  }
  column "internal_id" {
    null = true
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "attachment_disable" {
    null = false
    type = boolean
  }
  column "text_disable" {
    null = false
    type = boolean
  }
  column "title_template" {
    null = false
    type = text
  }
  column "body_template" {
    null = false
    type = text
  }
  column "kind" {
    null = false
    type = text
  }
  column "config" {
    null = false
    type = json
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  column "created_at" {
    null    = false
    type    = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "endpoints_internal_id_idx" {
    unique  = true
    columns = [column.internal_id]
  }
}
table "rules" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "internal" {
    null = false
    type = boolean
  }
  column "internal_id" {
    null = true
    type = text
  }
  column "name" {
    null = false
    type = text
  }
  column "expression" {
    null = false
    type = text
  }
  column "enable" {
    null = false
    type = boolean
  }
  column "updated_at" {
    null = false
    type = datetime
  }
  column "created_at" {
    null    = false
    type    = datetime
  }
  primary_key {
    columns = [column.id]
  }
  index "rules_internal_id_idx" {
    unique  = true
    columns = [column.internal_id]
  }
}
table "rules_to_endpoints" {
  schema = schema.main
  column "internal" {
    null = false
    type = boolean
  }
  column "rule_id" {
    null = false
    type = integer
  }
  column "endpoint_id" {
    null = false
    type = integer
  }
  foreign_key "endpoint_id" {
    columns     = [column.endpoint_id]
    ref_columns = [table.endpoints.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  foreign_key "rule_id" {
    columns     = [column.rule_id]
    ref_columns = [table.rules.column.id]
    on_update   = CASCADE
    on_delete   = CASCADE
  }
  index "rules_to_endpoints_rule_id_endpoint_id_idx" {
    columns = [column.rule_id, column.endpoint_id]
    unique  = true
  }
}
table "traces" {
  schema = schema.main
  column "id" {
    null = false
    type = integer
    auto_increment = true
  }
  column "request_id" {
    null = false
    type = text
  }
  column "source" {
    null = false
    type = text
  }
  column "seq" {
    null = false
    type = integer
  }
  column "action" {
    null = false
    type = text
  }
  column "level" {
    null = false
    type = text
  }
  column "data" {
    null = false
    type = json
  }
  column "created_at" {
    null    = false
    type    = datetime
  }
  primary_key {
    columns = [column.id]
  }
}
table "mailman_queue" {
  schema = schema.main
  column "message_id" {
    null = false
    type = integer
  }
  column "created_at" {
    null    = false
    type    = datetime
  }
  index "mailman_queue_message_id_idx" {
    columns = [column.message_id]
    unique  = true
  }
  foreign_key "message_id" {
    columns     = [column.message_id]
    ref_columns = [table.messages.column.id]
    on_update   = CASCADE 
    on_delete   = CASCADE
  }
}
schema "main" {
}
