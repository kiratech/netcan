BINARY_NAME=netcan
NETCAN_IMAGE_FQDN=docker.io/kiratech/netcan
TOOLS_IMAGE_FQDN=docker.io/kiratech/netcan-tools
PROJECT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CURDIR=
DEVIMAGEOPTS=-v $(PROJECT_DIR):/go/src/github.com/kiratech/netcan -w /go/src/github.com/kiratech/netcan $(TOOLS_IMAGE_FQDN):latest
BUILDCONTAINER=docker run --rm $(DEVIMAGEOPTS)
LDFLAGS='-extldflags "-static"'
TESTABLE_PACKAGES=./network ./pkg/mountinfo

ifeq ($(NOCONTAINER), 1)
	BUILDCONTAINER=
	CURDIR=$(PROJECT_DIR)/
endif

# We haven't released yet, for now the version is just the commit hash
VERSION=$(shell git rev-parse --verify HEAD)

.PHONY: docker-image
docker-image: build
	docker image build -t $(NETCAN_IMAGE_FQDN):$(VERSION) -f Dockerfile .
	docker image tag $(NETCAN_IMAGE_FQDN):$(VERSION) $(NETCAN_IMAGE_FQDN):latest

.PHONY: push-images
push-images: docker-image
	docker image push $(NETCAN_IMAGE_FQDN):$(VERSION)
	docker image push $(NETCAN_IMAGE_FQDN):latest
	docker image push $(TOOLS_IMAGE_FQDN):latest

.PHONY: build
build: $(BINARY_NAME)

$(BINARY_NAME): tools-image
	$(BUILDCONTAINER) sh -c "CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) ."

.PHONY: tools-image
tools-image:
	docker image build -t $(TOOLS_IMAGE_FQDN):latest -f Dockerfile.tools .

test:
	go test $(TESTABLE_PACKAGES)
