ARG ALPINE_VERSION=3.15



FROM alpine:${ALPINE_VERSION} AS build-env

WORKDIR /build

ARG GO_VERSION=1.17.2
ARG PROTOC_VERSION=3.18.1
ARG PROTOC_GEN_GO_VERSION=1.27.1
ARG PROTOC_GEN_GO_GRPC_VERSION=1.2.0

ARG BUILD_OS=linux
ARG BUILD_ARCH=amd64
ARG TARGET_OS=linux
ARG TARGET_ARCH=amd64

ENV GOPATH=/go
ENV GOCACHE=/go/cache
ENV GOOS="${TARGET_OS}"
ENV GOARCH="${TARGET_ARCH}"
ENV GOFLAGS="-v -mod=readonly -mod=vendor"
ENV GO111MODULE=on
ENV CGO_ENABLED=0

ENV PATH="${GOPATH}/bin:/usr/local/go/bin/:${PATH}"

RUN --mount=type=cache,target=/var/cache/apt set -eux; \
  apk add \
    git \
    libc6-compat

RUN set -eux; \
  wget -O go.tar.gz "https://dl.google.com/go/go${GO_VERSION}.${BUILD_OS}-${BUILD_ARCH}.tar.gz" ; \
  tar -C /usr/local/ -xzf go.tar.gz; \
  rm -f go.tar.gz

ENV PATH=/usr/local/bin:${PATH}

RUN go version



FROM build-env AS build

COPY . ./

RUN --mount=type=cache,target=/go/cache set -eux; \
  go vet ./...

RUN --mount=type=cache,target=/go/cache set -eux; \
  go test -ldflags "-s -w -extldflags '-static -fno-PIC'" ./...



FROM scratch AS end-of-file

RUN false
