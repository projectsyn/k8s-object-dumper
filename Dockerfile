FROM debian:10.5-slim as base

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
     bash \
     jq \
     less \
     moreutils \
     procps \
  && apt-get clean

FROM base as downloader

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
     ca-certificates \
     curl \
  && apt-get clean

RUN curl -sLo /tmp/kubectl "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
  && chmod +x /tmp/kubectl

RUN curl -sLo /tmp/krossa.tar.gz https://github.com/appuio/krossa/releases/download/v0.0.4/krossa_0.0.4_linux_amd64.tar.gz \
  && mkdir /tmp/krossa \
  && tar -xzf /tmp/krossa.tar.gz --directory /tmp/krossa \
  && ls /tmp/krossa

FROM base

RUN mkdir /data \
  && chown 1001:0 /data

COPY --from=downloader /tmp/kubectl /usr/local/bin
COPY --from=downloader /tmp/krossa/krossa /usr/local/bin
COPY dump-objects /usr/local/bin

USER 1001
ENTRYPOINT ["/usr/local/bin/dump-objects"]
CMD ["-d", "/data"]
