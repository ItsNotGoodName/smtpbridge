# smtpbridge

SMTP server that bridges email to other types of messaging services.

It only accepts attachments that are `image/png` or `image/jpeg`.

Do not expose this to the Internet, this is only intended to be used on your local network.

## Senders

List of example senders.

- Motion/AI detection from IP Cameras
- Notifications from Linux servers such as unattended updates

## Endpoints

The following message endpoints have been implemented.

- Telegram

## Configuration

Configuration file is located at `~/.smtpbridge.yml`.

```yaml
# Simple configuration
bridges:
  - name: test bridge
    endpoints:
      - test endpoint # Must match a name in the endpoints list

endpoints:
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
```

```yaml
# Full configuration
auth:
  enable: true # Enable auth, default 'false'
  username: user
  password: 12345678

database:
  type: bolt # Default ''
  db: /tmp/smtpbridge.db # Optional, default '$HOME/.smtpbridge/smtpbridge.db'
  attachments: /tmp/attachments # Optional, default '$HOME/.smtpbridge/attachments'

smtp:
  host: "localhost" # Host to listen on, default ''
  port: 1025 # Port to listen on, default 1025
  size: 26214400 # Max allowed size of email in bytes, default 26214400 (25 MB)

bridges:
  - name: test bridge
    filters:
      - to: foo@example.com # Filter based on to address
        from: bar@example.com # Filter based on from address
    only_text: false # When this is true, only the text of the email will be sent to endpoints
    only_attachments: false # When this is true, only the attachments of the email will be sent to endpoints
    endpoints:
      - test endpoint # Match a name in the endpoints list

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

## Todo

- SMTP authentication
- Regex from and to addresses
- Web interface
