#!/bin/sh

if [ -n "$APPRISE_ENABLE" ]; then
	apk add py3-pip

	if [ -n "$APPRISE_VERSION" ]; then
		pip3 install "apprise==$APPRISE_VERSION"
	else
		pip3 install apprise
	fi
fi

/usr/bin/smtpbridge --data-directory=/data
