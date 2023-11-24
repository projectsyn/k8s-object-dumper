FROM docker.io/debian:12.2-slim as base

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
     bash \
     jq \
     less \
     moreutils \
     procps \
  && rm -rf /var/lib/apt/lists/*

FROM base as downloader

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
     ca-certificates \
     curl \
  && rm -rf /var/lib/apt/lists/*

ARG K8S_VERSION=v1.18.20

RUN curl -sLo /tmp/kubectl "https://storage.googleapis.com/kubernetes-release/release/${K8S_VERSION}/bin/linux/amd64/kubectl" \
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
COPY must-exist /usr/local/share/k8s-object-dumper/
COPY known-to-fail /usr/local/share/k8s-object-dumper/

USER 1001
ENTRYPOINT ["/usr/local/bin/dump-objects"]
CMD ["-d", "/data"]
