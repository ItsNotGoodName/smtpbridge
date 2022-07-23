# smtpbridge

SMTP server that bridges email to other messaging services.

Do not expose this to the Internet, this is only intended to be used on a local network.

## Use Cases

- Receive pictures from IP cameras
- Receive system messages from Linux servers

## Configuration

Configuration file is located at `~/.smtpbridge.yml`.

```yaml
database: # Database
  type: memory # (memory)
  memory:
    limit: 100 # Max number of envelopes

storage: # Storage for attachment data
  type: memory # (memory)
  memory:
    limit: 30 # Max number of attachment data
    size: 104857600 # Max memory allocation, 100 MiB

http: # HTTP server
  enable: true
  host: ""
  port: 8080

smtp: # SMTP server
  enable: true
  host: ""
  port: 1025
  size: 26214400 # Max message size in bytes, 25 MiB
  auth: false
  username: ""
  password: ""

endpoints: # Endpoints for envelopes
  - name: telegram endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
  - name: console endpoint
    type: console
```

## Usage

```
smtpbridge
```

## To Do
