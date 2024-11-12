FROM docker.io/library/alpine:3.20 as runtime

RUN \
  apk add --update --no-cache \
    bash \
    curl \
    ca-certificates \
    tzdata

ENTRYPOINT ["k8s-object-dumper"]
COPY k8s-object-dumper /usr/bin/

USER 65536:0
