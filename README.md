# SMTPBridge

Bridge email to other messaging services.

**Do not expose this to the Internet, this is only intended to be used on a local network.**

## Use Cases

- Pictures from IP cameras
- System messages from servers

## Usage

```
smtpbridge
```

Run with database and storage in memory.

```
smtpbridge --memory
```

Restart when config file changes.

```
smtpbridge --watch
```

## Config

Config file is located at `~/.smtpbridge.yml`.

### Starter Config

This config prints emails received via SMTP to console.
The SMTP server listens on port `1025` and the HTTP server listens on port `8080`.
This saves emails to `~/.smtpbridge` directory.

```yaml
endpoints:
  - name: hello world
    type: console
```

### Full Config

```yaml
memory: false # Run with database and storage in memory

directory: ~/.smtpbridge # Default persistence directory

database: # Database
  type: bolt # (bolt, memory)
  memory:
    limit: 100 # Max number of envelopes

storage: # Storage
  type: file # (file, memory)
  memory:
    size: 104857600 # Max memory allocation in bytes, 100 MiB

http: # HTTP server
  disable: true # (false, true)
  host: ""
  port: 8080

smtp: # SMTP server
  disable: false # (false, true)
  host: ""
  port: 1025
  size: 26214400 # Max message size in bytes, 25 MiB
  username: ""
  password: ""

endpoints: # Endpoints for envelopes
  - name: example endpoint
    text_disable: false
    text_template: |
      FROM: {{ .Message.From }}
      SUBJECT: {{ .Message.Subject }}
      {{ .Message.Text }}
    attachments_disable: false
    type: console
  # Console
  - name: console endpoint
    type: console
  # Telegram
  - name: telegram endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111

bridges: # Bridges to endpoints, if this is empty then envelopes will always be sent to all endpoints
  # Send to 'console endpoint'
  - endpoints: console endpoint
  # Send to 'console endpoint' if the envelope is from 'foo@example.com' and is to 'bar@example.com'
  - from: foo@example.com
    to: bar@example.com
    endpoints: console endpoint
  # Send to all endpoints if the envelope is from 'foo@example.com' or 'baz@example.com'
  - filters:
      - from: foo@example.com
      - from: baz@example.com
  # Send to 'console endpoint' if the envelope to matches regex "@example\.com$"
  - to_regex: '@example\.com$'
    endpoints: console endpoint
  # Send to 'telegram endpoint' and 'console endpoint' if the envelope has more than 4 attachments
  - match_template: "{{ gt (len .Attachments) 4 }}"
    endpoints:
      - telegram endpoint
      - console endpoint
```

### Template

Each template has access to the [`Envelope`](./core/envelope/envelope.go) struct.
See [`text/template`](https://pkg.go.dev/text/template) on how to template.

## To Do

- Read mail files
- Save raw emails
- Windows installer
