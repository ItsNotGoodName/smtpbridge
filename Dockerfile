FROM alpine

ENV MIRROR=http://dl-cdn.alpinelinux.org/alpine

# Add support for TZ environment variable
RUN apk add --no-cache tzdata 

VOLUME /data

WORKDIR /config

ENTRYPOINT ["/usr/bin/smtpbridge", "--data-directory=/data"]

COPY smtpbridge /usr/bin/smtpbridge
