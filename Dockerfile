FROM alpine

ENV APPRISE_ENABLE=
ENV APPRISE_VERSION=

VOLUME /data
WORKDIR /config

ENTRYPOINT ["/entrypoint.sh"]

COPY entrypoint.sh /entrypoint.sh
COPY smtpbridge /usr/bin/smtpbridge
