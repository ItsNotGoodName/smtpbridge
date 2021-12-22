# smtpbridge

## Configuration

Configuration file at `~/.smtpbridge.yml`.

```yaml
bridges:
  - name: test bridge
    filters:
      - to: test@example.com
    endpoints:
      - test endpoint

endpoints:
  - name: test endpoint
    type: telegram
    config:
      token: 2222222222222222222222
      chat_id: 111111111111111111111
```
