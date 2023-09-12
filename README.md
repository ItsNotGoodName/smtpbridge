# SMTPBridge

[![GitHub](https://img.shields.io/github/license/itsnotgoodname/smtpbridge)](./LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge/tags)
[![GitHub last commit](https://img.shields.io/github/last-commit/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/itsnotgoodname/smtpbridge)](./go.mod)

Bridge email to other messaging services.

**Do not expose this to the Internet as this is only intended to be used on a local network.**

[![Screenshot](https://static.gurnain.com/github/smtpbridge/demo-small.png)](https://static.gurnain.com/github/smtpbridge/demo.png)

# Use Cases

- Pictures from IP cameras
- System messages from servers

# Usage

```
smtpbridge
```

## Show Version

```
smtpbridge -version
```

# Supported Endpoints

- console
- [Telegram](https://telegram.org/)
- [Shoutrrr](https://github.com/containrrr/shoutrrr)
- [Apprise](https://github.com/caronc/apprise)

## Apprise

Apprise requires Python to be installed along with the `apprise` package.

Install Apprise with the following command.

```
pip install apprise
```

If you are using a Python virtual environment, then set the `python_executable` config variable.

```yaml
python_executable: .venv/bin/python3
```

Make sure you install Apprise in that virtual environment.

# Config

Config file is loaded from one of the following locations in order.

- `config.yaml`
- `config.yml`
- `.smtpbridge.yaml`
- `.smtpbridge.yml`
- `~/.smtpbridge.yaml`
- `~/.smtpbridge.yml`
- `/etc/smtpbridge.yaml`
- `/etc/smtpbridge.yml`

## Simple Config

This config prints emails received via SMTP to console.

```yaml
endpoints:
  hello_world:
    kind: console

rules:
  hello_world:
```

## Full Config

```yaml
# Used for development
debug: false

# Used by HTTP, ...
time_format: 12h # (12h, 24h)

# Maximum message size for envelopes
max_payload_size: 25 MB

# Directory for storing data
data_directory: smtpbridge_data

# Python executable
python_executable: python3

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

  # Public URL
  url: "" # (http://127.0.0.1:8080)

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
      # https://containrrr.dev/shoutrrr/v0.8/services/overview/
      urls: telegram://token@telegram?chats=@channel-1[,chat-id-1,...]

  # Apprise
  shoutrrr_endpoint:
    kind: apprise
    config:
      # https://github.com/caronc/apprise#supported-notifications
      urls: tgram://bottoken/ChatID

  # Full example
  example_endpoint:
    kind: console
    name: Example Endpoint
    # Do not send title and body
    text_disable: false
    # Do not send attachments
    attachment_disable: false
    # Custom template for title
    title_template: {{ .Message.Subject }}
    # Custom template for body
    body_template: {{ .Message.Text }}

rules:
  example_rule:
    name: Example Rule
    expression: or (eq .Message.Subject "cam-1") (eq .Message.Subject "cam-2")
    endpoints:
      - console_endpoint
```

## Templates

See [`text/template`](https://pkg.go.dev/text/template) on how to template.

Each `*_template` config variable has access to the [Envelope](./internal/models/models.go) model via the `.` operator.

The following custom functions are available in endpoint templates.

| Name      | Description                             | Example                                                           |
| --------- | --------------------------------------- | ----------------------------------------------------------------- |
| PermaLink | Permanent HTTP link to the given model. | `{{ PermaLink .Message }}` => `http://127.0.0.1:8080/envelopes/1` |

## Expressions

Rule expressions are just [Templates](#templates) without `{{ }}` and without the custom functions.

# Docker

## Docker Compose

```yaml
version: "3"
services:
  smtpbridge:
    image: ghcr.io/itsnotgoodname/smtpbridge:latest
    container_name: smtpbridge
    environment:
      APPRISE_ENABLE: "true" # Optional
      APPRISE_VERSION: "1.5.0" # Optional
      SMTPBRIDGE_CONFIG_YAML: | # Optional
        endpoints:
          hello_world:
            kind: console

        rules:
          hello_world:
    ports:
      - 1025:1025
      - 8080:8080
    volumes:
      - /path/to/data:/data
      - /path/to/config:/config # Optional
      - /etc/timezone:/etc/timezone:ro # Optional
      - /etc/localtime:/etc/localtime:ro # Optional
    restart: unless-stopped
```

## Docker CLI

```sh
docker run -d \
  --name=smtpbridge \
  -e APPRISE_ENABLE=true `# Optional` \
  -e APPRISE_ENABLE=1.5.0 `# Optional` \
  -p 1025:1025 \
  -p 8080:8080 \
  -v /path/to/data:/data \
  -v /path/to/config:/config `# Optional` \
  -v /etc/timezone:/etc/timezone:ro `# Optional` \
  -v /etc/localtime:/etc/localtime:ro `# Optional` \
  --restart unless-stopped \
  ghcr.io/itsnotgoodname/smtpbridge:latest
```

# Development

The following programs are required.

- Make
- Go
- pnpm

## Make

**You should look at the [Makefile](./Makefile) before running any of the following commands.**

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

# To Do

- feat: read [mbox](https://access.redhat.com/articles/6167512) files
- feat: CRUD endpoints
- feat: IMAP for viewing mail
- feat: JSON API
- feat: Windows installer
- fix: chrome keeps thinking some HTTP pages are French
