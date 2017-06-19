BINARY_NAME=netcan
NETCAN_IMAGE_FQDN=docker.io/fntlnz/netcan:latest
TOOLS_IMAGE_FQDN=docker.io/fntlnz/netcan-tools:latest
PROJECT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CURDIR=
DEVIMAGEOPTS=-v $(PROJECT_DIR):/go/src/github.com/fntlnz/netcan -w /go/src/github.com/fntlnz/netcan $(TOOLS_IMAGE_FQDN)
BUILDCONTAINER=docker run --rm $(DEVIMAGEOPTS)
LDFLAGS='-extldflags "-static"'
TESTABLE_PACKAGES=./network ./proc

ifeq ($(NOCONTAINER), 1)
	BUILDCONTAINER=
	CURDIR=$(PROJECT_DIR)/
endif

.PHONY: docker-image
docker-image: build
	docker build -t $(NETCAN_IMAGE_FQDN) -f Dockerfile .

.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): tools-image
	$(BUILDCONTAINER) sh -c "CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) ."

.PHONY: tools-image
tools-image:
	docker build -t $(TOOLS_IMAGE_FQDN) -f Dockerfile.tools .

test:
	go test $(TESTABLE_PACKAGES)
