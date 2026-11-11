ARG BASE_IMAGE=alpine:3.22@sha256:9db2934886712d877c932e07249c75a59e6ca27adfbccc05261987a17bef2d63
FROM ${BASE_IMAGE}

ARG BUILDARCH
WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates
COPY ./${BUILDARCH}/release /app/
COPY LICENSE /app/LICENSE
VOLUME /app/data

EXPOSE 21114
CMD ["./apimain"]
