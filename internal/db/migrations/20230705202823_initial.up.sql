CREATE TABLE IF NOT EXISTS messages (
    id INTEGER PRIMARY KEY,
    from_ TEXT NOT NULL,
    to_ JSON NOT NULL,
    subject TEXT NOT NULL,
    text TEXT NOT NULL,
    html TEXT NOT NULL,
    date DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

--bun:split

CREATE TABLE IF NOT EXISTS attachments (
    id INTEGER PRIMARY KEY,
    message_id INTEGER,
    name TEXT NOT NULL,
    mime TEXT NOT NULL,
    extension TEXT NOT NULL,
    FOREIGN KEY (message_id) REFERENCES messages (id) ON DELETE SET NULL
);

--bun:split

CREATE TABLE IF NOT EXISTS endpoints (
    id INTEGER PRIMARY KEY,
    internal BOOLEAN NOT NULL,
    internal_id TEXT NOT NULL,
    name TEXT NOT NULL,
    attachment_disable BOOLEAN NOT NULL,
    text_disable BOOLEAN NOT NULL,
    body_template TEXT NOT NULL,
    kind TEXT NOT NULL,
    config JSON NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (internal_id)
);

--bun:split

CREATE TABLE IF NOT EXISTS rules (
    id INTEGER PRIMARY KEY,
    internal BOOLEAN NOT NULL,
    internal_id TEXT NOT NULL,
    name TEXT NOT NULL,
    expression TEXT NOT NULL,
    enable BOOLEAN NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (internal_id)
);

--bun:split

CREATE TABLE rules_to_endpoints (
    internal BOOLEAN NOT NULL,
    rule_id INTEGER NOT NULL,
    endpoint_id TEXT NOT NULL,
    updated_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (rule_id) REFERENCES rules (id) ON DELETE CASCADE,
    FOREIGN KEY (endpoint_id) REFERENCES endpoints (id) ON DELETE CASCADE,
    UNIQUE (rule_id, endpoint_id)
);
