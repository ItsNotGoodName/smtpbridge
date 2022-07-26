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

## Configuration

Configuration file is located at `~/.smtpbridge.yml`.

### Starter Configuration

This config prints emails received via SMTP to console.
The SMTP server listens on port `1025` and the HTTP server listens on port `8080`.
This saves emails to `~/.smtpbridge` directory.
The `database` and `storage` keys can be removed to run this entirely in memory.

```yaml
database:
  type: bolt

storage:
  type: directory

endpoints:
  - name: hello world
    type: console
```

### Full Configuration

```yaml
directory: ~/.smtpbridge # Default persistence directory

database: # Database
  type: memory # (memory, bolt)
  memory:
    limit: 100 # Max number of envelopes

storage: # Storage for attachment's data
  type: memory # (memory, directory)
  memory:
    size: 104857600 # Max memory allocation in bytes, 100 MiB

http: # HTTP server
  enable: true # (true, false)
  host: ""
  port: 8080

smtp: # SMTP server
  enable: true # (true, false)
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

## To Do

- Read mail files
- Windows installer
