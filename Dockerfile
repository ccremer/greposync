FROM docker.io/library/alpine:3.14 as runtime

ENTRYPOINT ["gsync"]

RUN \
  apk add --no-cache curl bash git openssh

COPY gsync /usr/bin/
USER 1000:0
