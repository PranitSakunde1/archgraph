BINARY ?= archgraph
IMAGE ?= archgraph:dev
VERSION ?= dev
COMMIT ?= unknown
BUILD_DATE ?= unknown

.PHONY: test build docker-build docker-scan clean

test:
	go test ./...

build:
	CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT) -X main.buildDate=$(BUILD_DATE)" -o $(BINARY) ./cmd/archgraph

docker-build:
	docker build --pull --build-arg VERSION=$(VERSION) --build-arg COMMIT=$(COMMIT) --build-arg BUILD_DATE=$(BUILD_DATE) -t $(IMAGE) .

docker-scan:
	@test -n "$(WORKSPACE)" || (echo "Set WORKSPACE=/absolute/path/to/repos" && exit 1)
	@mkdir -p "$(OUTPUT)"
	docker run --rm -v "$(WORKSPACE):/workspace:ro" -v "$(OUTPUT):/output" $(IMAGE) scan --path /workspace --output /output

clean:
	rm -f $(BINARY) archgraph.exe
