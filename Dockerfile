FROM alpine

VOLUME /data

WORKDIR /config

RUN touch /config/config.yaml

ENTRYPOINT ["/usr/bin/smtpbridge", "--data-directory=/data"]

COPY smtpbridge /usr/bin/smtpbridge
