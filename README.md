# smtpbridge

## Configuration

Configuration file at `~/.smtpbridge.yml`.

```yaml
smtp:
  port: 1025
  size: 26214400 # 25 MB

bridges:
  - name: test bridge
    email_to: foo@example.com
    email_from: bar@example.com
    only_text: false
    only_attachments: false
    endpoints:
      - test endpoint

endpoints:
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
```
