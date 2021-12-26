# smtpbridge

SMTP server that bridges email to other types of messaging services.

## Message Endpoints

The following message endpoints have been implemented.

- Telegram

## Configuration

Configuration file at `~/.smtpbridge.yml`.

```yaml
smtp:
  host: "" # Host to listen on
  port: 1025 # Port to listen on
  size: 26214400 # 25 MB, max allowed size of email in bytes

bridges:
  - name: test bridge
    email_to: foo@example.com
    email_from: bar@example.com
    only_text: false # When this is true, only the text of the email will be sent to endpoints
    only_attachments: false # When this is true, only the attachments of the email will be sent to endpoints
    endpoints:
      - test endpoint # Must match a name in the endpoints list

endpoints:
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
```

## Usage

```
smtpbridge server
```
