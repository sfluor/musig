BINARY = musig
GOARCH = amd64

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${BRANCH}:${COMMIT}"

# Enable go modules
GOCMD = GO111MODULE=on go

# Build the project
all: build

.PHONY: build
build:
	${GOCMD} build ${LDFLAGS} -o ./bin/${BINARY} ./cmd/musig/main.go

.PHONY: linux
linux:
	GOOS=linux GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} .

.PHONY: macos
macos:
	GOOS=darwin GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-macos-${GOARCH} .

.PHONY: windows
windows:
	GOOS=windows GOARCH=${GOARCH} ${GOCMD} build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe .

cross: linux macos windows

.PHONY: test
test:
	${GOCMD} get -v ./...; \
	${GOCMD} vet $$(go list ./... | grep -v /vendor/); \
	${GOCMD} test -v -race ./...; \

.PHONY: fmt
fmt:
	${GOCMD} fmt $$(go list ./... | grep -v /vendor/)

.PHONY: tidy
tidy:
	${GOCMD} mod tidy

.PHONY: download
download:
	./scripts/dl_dataset.sh
