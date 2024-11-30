# SMTPBridge

[![GitHub](https://img.shields.io/github/license/itsnotgoodname/smtpbridge)](./LICENSE)
[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge/tags)
[![GitHub last commit](https://img.shields.io/github/last-commit/itsnotgoodname/smtpbridge)](https://github.com/ItsNotGoodName/smtpbridge)
[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/itsnotgoodname/smtpbridge)](./go.mod)

Bridge email to other messaging services.

**Do not expose this to the Internet as this is only intended to be used on a local network.**

[![Screenshot](https://static.gurnain.com/github/smtpbridge/demo-small.png)](https://static.gurnain.com/github/smtpbridge/demo.png)

# Features

- Receive email from SMTP or HTTP as envelopes
- Send envelopes to [endpoints](#supported-endpoints) with [templates](#templates)
- Create [rules](#expressions) for matching envelopes with endpoints
- View and manage application through the Web UI
- Delete stale envelopes with a retention policy
- Monitor application with healthcheck

# Use Cases

- Pictures from IP cameras (e.g. AI Tripwire, ...)
- System messages from servers and applications (e.g. Debian, Nextcloud, UniFi Network Application, ...)

# Usage

```
smtpbridge
```

# Supported Endpoints

- Console
- [Telegram](https://telegram.org/)
- [Shoutrrr](https://github.com/containrrr/shoutrrr)
- [Apprise](https://github.com/caronc/apprise)
- [Script](#script)

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

## Script

This allows you to run an arbitrary script as an endpoint.
stdin is a JSON encoded envelope such as the following.

```json
{
  "title": "Test Subject",
  "body": "Test Body",
  "attachments": [
    {
      "path": "http://127.0.0.1:8080/apple-touch-icon.png",
      "name": "Test Attachment"
    }
  ]
}
```

The `path` of an attachment can be a URL or a file path.

Please note that scripts runs concurrently.

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

CLI flags take priority over config files and environment.

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
time_format: 12h # [12h, 24h]

# Directory for storing data
data_directory: smtpbridge_data

# Python executable used by Apprise
python_executable: python3

# Healthcheck enables verification that the application has not crashed or lost network access
# You can use a third party service such as healthchecks.io
healthcheck:
  # URL to fetch, empty means health checking is disabled
  url: "" # [https://hc-ping.com/cb8bcf81-d3c4-4c98-85a6-734c3b7ddb2b, ...]

  # Interval between each fetch
  interval: 5m # [5m, 5h45m, ...]

  # Run on startup
  startup: false

# Mailman handles sending envelopes to the configured endpoints
mailman:
  # Number of concurrent workers
  workers: 1

# Retention policy will delete resources that pass the configured policy
retention:
  # Envelopes in database
  envelope_count: # [0, 100, 250, ...]
  envelope_age: # [5m, 5h45m, ...]

  # Attachment files in file store
  attachment_size: # [100 MB, 1 GB, ...]

  # Traces in database
  trace_age: 168h # 7 days [5m, 5h45m, ...]

# HTTP server
http:
  disable: false
  host: "" # [127.0.0.1, ...]
  port: 8080

  # Authentication is disabled if both username and password are empty
  username: ""
  password: ""

  # Public URL used for creating links
  url: "" # [http://127.0.0.1:8080, ...]

# SMTP server
smtp:
  disable: false
  host: "" # [127.0.0.1, ...]
  port: 1025

  # Authentication is disabled if both username and password are empty
  username: ""
  password: ""

  # Maximum payload size
  max_payload_size: 25 MB # [100 MB, 1 GB, ...]

# Endpoints for envelopes
endpoints:
  # Console
  console_endpoint:
    kind: console

  # Telegram
  telegram_endpoint:
    kind: telegram
    config:
      # https://core.telegram.org/bots/features#creating-a-new-bot
      token: 2222222222222222222222
      # https://stackoverflow.com/a/32572159
      chat_id: 111111111111111111111

  # Shoutrrr (can only send text)
  shoutrrr_endpoint:
    kind: shoutrrr
    config:
      # https://containrrr.dev/shoutrrr/v0.8/services/overview/
      urls: telegram://token@telegram?chats=@channel-1[,chat-id-1,...]

  # Apprise
  apprise_endpoint:
    kind: apprise
    config:
      # https://github.com/caronc/apprise#supported-notifications
      urls: tgram://bottoken/ChatID

  # Script
  script_endpoint:
    kind: script
    config:
      file: my-script.py

  # Full example
  example_endpoint:
    kind: console
    name: Example Endpoint
    # Do not send title and body
    text_disable: false
    # Do not send attachments
    attachment_disable: false
    # Custom template for title
    title_template: "{{ .Message.Subject }}"
    # Custom template for body
    body_template: "{{ .Message.Text }}"

rules:
  example_rule:
    name: Example Rule
    expression: or (eq .Message.Subject "cam-1") (eq .Message.Subject "cam-2")
    endpoints:
      - console_endpoint
```

## Templates

See [`text/template`](https://pkg.go.dev/text/template) on how to template.

Each `*_template` has access to the [Envelope](./internal/models/envelope.go) model via the `.` operator.

The following custom functions are available in endpoint templates.

| Name      | Description                              | Example                                                           |
| --------- | ---------------------------------------- | ----------------------------------------------------------------- |
| PermaLink | Permanent HTTP link for the given model. | `{{ PermaLink .Message }}` => `http://127.0.0.1:8080/envelopes/1` |

## Expressions

Rule expressions are just [`text/template`](https://pkg.go.dev/text/template) without `{{ }}`.
The [Envelope](./internal/models/envelope.go) model can be accessed via the `.` operator.
They should always evaluate to a boolean expression.

Example:

```
or
  (eq .Message.Subject "cam-1")
  (.Message.To.EQ "my-name@example.com")
  (eq .Message.From "unleashed@example.com")
```

This expression will pass if one of the following is true.

- The message's subject equals "cam-1"
- The message is to "<my-name@example.com>"
- The message is from "<unleashed@example.com>"

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

Install tooling.

```
make tooling
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

- refactor: WAY TOO MANY TOOLS TO BUILD THE PROGRAM, REMOVE SOME
- feat: read [mbox](https://access.redhat.com/articles/6167512) files
- feat: IMAP for viewing mail
- feat: OpenAPI
- feat: Windows installer
- fix: chrome keeps thinking some HTTP pages are French
