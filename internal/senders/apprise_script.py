import sys, json
import apprise

data = json.load(sys.stdin)

apobj = apprise.Apprise()

for url in data["urls"]:
    apobj.add(url)

paths = []
for attachment in data['attachments']:
    paths.append(attachment['path'])

attach = apprise.AppriseAttachment(paths=paths)

with apprise.LogCapture(level=apprise.logging.INFO) as logs:
    if not apobj.notify(body=data['body'], title=data['title'], attach=attach):
        print(logs.getvalue(), file=sys.stderr) # type: ignore
        sys.exit(1)
