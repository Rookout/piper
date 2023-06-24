FROM golang:1.20-alpine3.16 as builder

WORKDIR /piper

RUN apk update && apk add --no-cache \
    git \
    make \
    wget \
    curl \
    gcc \
    bash \
    ca-certificates \
    musl-dev \
    zlib-static \
    build-base

COPY go.mod .
COPY go.sum .
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go mod tidy

RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build  go build -gcflags='all=-N -l' -tags=alpine -buildvcs=false  -trimpath ./cmd/piper


FROM alpine:3.16 as piper-release

ENV GIN_MODE=release

USER 1001

COPY /.git /.git

COPY --chown=1001 --from=builder /piper/piper /bin

ENTRYPOINT [ "piper" ]