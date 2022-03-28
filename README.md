# smtpbridge

SMTP server that bridges email to other types of messaging services.

It only accepts attachments that are `image/png` or `image/jpeg`.

Do not expose this to the Internet, this is only intended to be used on your local network.

## Use Cases

- Receive motion/AI detection from IP cameras
- Receive notifications from Linux servers such as unattended updates

## Endpoints

The following endpoints for messages have been implemented.

- Telegram

## Configuration

Configuration file is located at `~/.smtpbridge.yml`.

```yaml
# Simple configuration
bridges:
  - name: test bridge
    endpoints:
      - name: test endpoint # Match a name in the endpoints list

endpoints:
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
```

```yaml
# Full configuration
storage:
  path: /tmp/smtpbridge # Path to storage location, Default '$HOME/.smtpbridge'
  size: 2147483648 #  Max size of storage location in bytes, Default 2147483648 (2 GiB)

database:
  type: bolt # Database type, Options ('', 'bolt')

http: # HTTP server that shows past messages
  enable: true # Enable http server, Default 'false'
  host: "127.0.0.1" # Host to listen on, Default ''
  port: 9000 # Port to listen on, Default 8080

smtp: # SMTP server that receives emails
  host: "127.0.0.1" # Host to listen on, Default ''
  port: 1025 # Port to listen on, Default 1025
  size: 26214400 # Max allowed size of email in bytes, Default 26214400 (25 MiB)
  auth: true # Enable auth, Default 'false'
  username: user # Default ''
  password: 12345678 # Default ''

bridges: # Bridges modify and check if messages should be sent to endpoints
  - name: test bridge
    no_text: false # When this is true, text will not be sent to endpoints
    no_attachments: false # When this is true, attachments will not be sent to endpoints
    min_attachments: 1 # Message must have at least this amount of attachments, Default 0
    filters:
      - to: foo@example.com # Filter based on to address
        to_regex: "foo" # To regex takes priority over to, must be surrounded by quotation marks
        from: bar@example.com # Filter based on from address
        from_regex: "bar" # From regex takes priority over from, must be surrounded by quotation marks
    endpoints:
      - name: test endpoint # Match a name in the endpoints list
        no_text: false # When this is true, text will not be sent to endpoints
        no_attachments: false # When this is true, attachments will not be sent to endpoints

endpoints: # Endpoints send messages to messaging services such as Telegram
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
  - name: mock endpoint
    type: mock # This endpoints prints out the message to console
```

## Usage

```
smtpbridge server
```

## Todo
