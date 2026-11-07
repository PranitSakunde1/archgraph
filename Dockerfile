# syntax=docker/dockerfile:1.7
# Version-pinned build inputs. CI should additionally pin both FROM references
# to approved OCI digests for bit-for-bit rebuilds.
ARG GO_VERSION=1.23.2
ARG ALPINE_VERSION=3.21.3
ARG PYTHON_IMAGE=python:3.12.8-slim-bookworm

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

FROM ${PYTHON_IMAGE} AS runtime
ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PIP_DISABLE_PIP_VERSION_CHECK=1 \
    HOME=/tmp
COPY requirements-graphify.txt /tmp/requirements-graphify.txt
RUN python -m pip install --no-cache-dir --no-compile -r /tmp/requirements-graphify.txt \
    && command -v graphify \
    && groupadd --gid 10001 archgraph \
    && useradd --uid 10001 --gid 10001 --no-create-home --shell /usr/sbin/nologin archgraph \
    && mkdir -p /workspace /output \
    && chown -R archgraph:archgraph /workspace /output
COPY --from=build /out/archgraph /usr/local/bin/archgraph
WORKDIR /workspace
USER 10001:10001
ENTRYPOINT ["/usr/local/bin/archgraph"]
CMD ["scan", "--path", "/workspace", "--output", "/output"]
