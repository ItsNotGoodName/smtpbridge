# SMTPBridge

[![GitHub](https://img.shields.io/github/license/itsnotgoodname/smtpbridge)](./LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge/tags)
[![GitHub last commit](https://img.shields.io/github/last-commit/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/itsnotgoodname/smtpbridge)](./go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ItsNotGoodName/smtpbridge)](https://goreportcard.com/report/github.com/ItsNotGoodName/smtpbridge)

Bridge email to other messaging services.

**Do not expose this to the Internet, this is only intended to be used on a local network.**

## Use Cases

- Pictures from IP cameras
- System messages from servers

## Usage

```
smtpbridge
```

### Show Version

```
smtpbridge --version
```

## Supported Endpoints

- Console
- [Telegram](https://telegram.org/)
- [Shoutrrr](https://github.com/containrrr/shoutrrr)

## Config

Config file is loaded from one of the following locations.

- `config.yaml`
- `config.yml`
- `.smtpbridge.yaml`
- `.smtpbridge.yml`
- `~/.smtpbridge.yaml`
- `~/.smtpbridge.yml`
- `/etc/smtpbridge.yaml`
- `/etc/smtpbridge.yml`

### Simple Config

This config prints emails received via SMTP to console.

```yaml
endpoints:
  hello_world:
    kind: console

rules:
  hello_world:
```

### Full Config

```yaml
# Maximum message size for envelopes
max_payload_size: 25 MB

# Directory for storing data
data_directory: smtpbridge_data

# Retention policy for envelopes and attachment files
retention:
  # # Retention policy for envelopes in database
  # envelope_count: # (0, 100, 250, ...)
  # envelope_age: # (5m, 5h45m, ...)

  # # Retention policy for attachment files on disk
  # attachment_size: # (100 MB, 1 GB, ...)

# HTTP server
http:
  disable: false # (false, true)
  host: ""
  port: 8080

  # Authentication is disabled if both username and password are empty
  username: ""
  password: ""

# SMTP server
smtp:
  disable: false # (false, true)
  host: ""
  port: 1025

  # Authentication is disabled if both username and password are empty
  username: ""
  password: ""

# Endpoints for envelopes
endpoints:
  # Full example
  example_endpoint:
    kind: console
    # Do not send any text
    text_disable: false
    # Custom template for body
    body_template: |
      {{ .Message.Text }}
    # Do not send any attachments
    attachment_disable: false

  # Console
  console_endpoint:
    kind: console

  # Telegram
  telegram_endpoint:
    kind: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111

  # Shoutrrr (can only send text)
  shoutrrr_endpoint:
    kind: shoutrrr
    config:
      # https://containrrr.dev/shoutrrr/0.7/services/overview/
      urls: |
        bark://devicekey@host
        discord://token@id
        smtp://username:password@host:port/?from=fromAddress&to=recipient1[,recipient2,...]
        gotify://gotify-host/token
        googlechat://chat.googleapis.com/v1/spaces/FOO/messages?key=bar&token=baz
        ifttt://key/?events=event1[,event2,...]&value1=value1&value2=value2&value3=value3
        join://shoutrrr:api-key@join/?devices=device1[,device2, ...][&icon=icon][&title=title]
        mattermost://[username@]mattermost-host/token[/channel]
        matrix://username:password@host:port/[?rooms=!roomID1[,roomAlias2]]
        ntfy://username:password@ntfy.sh/topic
        opsgenie://host/token?responders=responder1[,responder2]
        pushbullet://api-token[/device/#channel/email]
        pushover://shoutrrr:apiToken@userKey/?devices=device1[,device2, ...]
        rocketchat://[username@]rocketchat-host/token[/channel|@recipient]
        slack://[botname@]token-a/token-b/token-c
        teams://group@tenant/altId/groupOwner?host=organization.webhook.office.com
        telegram://token@telegram?chats=@channel-1[,chat-id-1,...]
        zulip://bot-mail:bot-key@zulip-domain/?stream=name-or-id&topic=name

rules:
  example_rule:
    name: Example Rule
    expression: or (eq .Message.Subject "cam-1") (eq .Message.Subject "cam-2")
    endpoints:
      - console_endpoint
```

### Template

Each template has access to [`Envelope`](./internal/envelope/envelope.go) via the `.` operator.

See [`text/template`](https://pkg.go.dev/text/template) on how to template.

## Docker

### Docker Compose

```yaml
version: "3"
services:
  smtpbridge:
    image: ghcr.io/itsnotgoodname/smtpbridge:latest
    container_name: smtpbridge
    environment:
      SMTPBRIDGE_CONFIG_YAML: | # Config by embedding YAML
        endpoints:
          hello_world:
            kind: console

        rules:
          hello_world:
    ports:
      - 1025:1025
      - 8080:8080
    volumes:
      - /path/to/config:/config # Config by creating config.yaml file in /config
      - /path/to/data:/data
      - /etc/timezone:/etc/timezone:ro
      - /etc/localtime:/etc/localtime:ro
    restart: unless-stopped
```

### Docker CLI

```sh
docker run -d \
  --name=smtpbridge \
  -p 1025:1025 \
  -p 8080:8080 \
  -v /path/to/config:/config \
  -v /path/to/data:/data \
  -v /etc/timezone:/etc/timezone:ro \
  -v /etc/localtime:/etc/localtime:ro \
  --restart unless-stopped \
  ghcr.io/itsnotgoodname/smtpbridge:latest
```

## Development

The following programs are required.

- Make
- Go
- pnpm

Install dependencies.

```
make dep
```

Start the Go server.

```
make dev
```

Start Vite.

```
make dev-web
```

## To Do

- Reload on config change
- Add [Apprise](https://github.com/caronc/apprise) endpoint
- HTTP and SMTP auth
- CRUD endpoints and rules
- SQLite full text search
- Read mailbox files
- Save raw emails
- JSON API
- Windows installer
