FROM alpine
VOLUME /data
WORKDIR /config
ENTRYPOINT ["/usr/bin/smtpbridge", "--data-directory=/data"]
COPY smtpbridge /usr/bin/smtpbridge
