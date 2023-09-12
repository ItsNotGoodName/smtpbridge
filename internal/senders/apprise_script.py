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

apobj.notify(
    body=data['body'],
    title=data['title'],
    attach=attach
)
