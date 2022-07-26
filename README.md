# SMTPBridge

SMTP server that bridges email to other messaging services.

Do not expose this to the Internet, this is only intended to be used on a local network.

## Use Cases

- Receive pictures from IP cameras
- Receive system messages from Linux servers

## Configuration

Configuration file is located at `~/.smtpbridge.yml`.

```yaml
directory: ~/.smtpbridge # Default persistence directory

database: # Database
  type: memory # (memory)
  memory:
    limit: 100 # Max number of envelopes

storage: # Storage for attachment data
  type: memory # (memory, directory)
  memory:
    size: 104857600 # Max memory allocation, 100 MiB

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
    type: console
    text_disable: false
    text_template: |
      FROM: {{ .Message.From }}
      SUBJECT: {{ .Message.Subject }}
      {{ .Message.Text }}
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
  - endpoints:
      - console endpoint
  # Send to all endpoints if the envelope is from 'example@example.com'
  - from: example@example.com
  # Send to 'console endpoint' if the envelope is from 'example@example.com' and is to 'test@example.com'
  - from: example@example.com
    to: test@example.com
    endpoints:
      - console endpoint
  # Send to 'console endpoint' if the envelope to matches regex "@example\.com$"
  - to_regex: "@example\.com$"
    endpoints:
      - console endpoint
  # Send to 'telegram endpoint' if the envelope has more than 4 attachments
  - match_template: "{{ gt (len .Attachments) 4 }}"
    endpoints:
      - telegram endpoint
```

## Usage

```
smtpbridge
```

## To Do
