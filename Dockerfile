ARG BASE_IMAGE=ubuntu:20.04
ARG ALPINE_IMAGE=docker.io/library/alpine:3.19.0
ARG BUILDER_IMAGE=golang:1.22

FROM ${BUILDER_IMAGE} AS builder

ARG TARGETOS
ARG TARGETARCH
ARG NOOPT
ARG NOSTRIP
ARG VERSION
ARG GIT_COMMIT

WORKDIR /go/src/bucket-proxy
RUN --mount=type=bind,readwrite,target=/go/src/bucket-proxy --mount=target=/root/.cache,type=cache --mount=target=/go/pkg,type=cache \
  make GOARCH=${TARGETARCH} GOOS=${TARGETOS} \
  VERSION=${VERSION} GIT_COMMIT=${GIT_COMMIT} NOOPT=${NOOPT} NOSTRIP=${NOSTRIP} \
  DESTDIR=/tmp/install/${TARGETOS}/${TARGETARCH} install

FROM ${ALPINE_IMAGE} AS certs
RUN apk --update add ca-certificates

FROM ${BASE_IMAGE} AS release
ARG TARGETOS
ARG TARGETARCH

LABEL maintainer="vietanhs0817@gmail.com"
LABEL org.opencontainers.image.source=https://github.com/vietanhduong/bucket-proxy
LABEL org.opencontainers.image.description="A simple bucke proxy server"
LABEL org.opencontainers.image.licenses=MIT

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /tmp/install/${TARGETOS}/${TARGETARCH}/usr/bin/* /usr/bin/
