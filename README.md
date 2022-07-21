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

storage: # Storage for attachments
  type: memory # (memory)

http: # HTTP server
  enable: true
  host: ""
  port: 8080

smtp: # SMTP server
  enable: true
  host: ""
  port: 1025
  size: 26214400 # Max message size in bytes (25 MiB)
  auth: false
  username: ""
  password: ""
```

## Usage

```
smtpbridge
```

## Todo
