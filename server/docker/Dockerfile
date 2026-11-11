ARG API_IMAGE=lejianwen/rustdesk-api:latest
FROM ${API_IMAGE}

ARG S6_OVERLAY_VERSION=3.2.0.0
ARG S6_ARCH=x86_64
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz /tmp
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${S6_ARCH}.tar.xz /tmp
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-noarch.tar.xz.sha256 /tmp
ADD https://github.com/just-containers/s6-overlay/releases/download/v${S6_OVERLAY_VERSION}/s6-overlay-${S6_ARCH}.tar.xz.sha256 /tmp
RUN \
  cd /tmp && \
  sha256sum -c s6-overlay-noarch.tar.xz.sha256 && \
  sha256sum -c s6-overlay-${S6_ARCH}.tar.xz.sha256 && \
  tar -C / -Jxpf /tmp/s6-overlay-noarch.tar.xz && \
  tar -C / -Jxpf /tmp/s6-overlay-${S6_ARCH}.tar.xz && \
  rm /tmp/s6-overlay*.tar.xz /tmp/s6-overlay*.sha256 && \
  if [ ! -e /var/run ]; then ln -s /run /var/run; fi

COPY rootfs /
RUN chmod 0755 /usr/bin/healthcheck.sh \
    /etc/s6-overlay/s6-rc.d/*/run \
    /etc/s6-overlay/s6-rc.d/key-secret/up.real

ENV RELAY=relay.example.com
ENV ENCRYPTED_ONLY=0

EXPOSE 21114 21115 21116 21116/udp 21117 21118 21119

HEALTHCHECK --interval=10s --timeout=5s CMD /usr/bin/healthcheck.sh

WORKDIR /app

VOLUME /app/data

VOLUME /data

ENTRYPOINT ["/init"]
