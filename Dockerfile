# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.23.2
ARG ALPINE_VERSION=3.21.3

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY cmd ./cmd
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH:-amd64} \
    go build -trimpath \
      -ldflags="-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.buildDate=${BUILD_DATE}" \
      -o /out/archgraph ./cmd/archgraph

FROM alpine:${ALPINE_VERSION} AS runtime
RUN apk add --no-cache ca-certificates \
    && addgroup -S -g 10001 archgraph \
    && adduser -S -D -H -u 10001 -G archgraph archgraph \
    && mkdir -p /workspace /output \
    && chown -R archgraph:archgraph /workspace /output
COPY --from=build /out/archgraph /usr/local/bin/archgraph
ENV HOME=/tmp
WORKDIR /workspace
USER 10001:10001
ENTRYPOINT ["/usr/local/bin/archgraph"]
CMD ["scan", "--path", "/workspace", "--output", "/output"]
