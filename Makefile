.PHONY: build tools-image

TOOLS_IMAGE_FQDN=docker.io/fntlnz/netcan-tools:latest
PROJECT_DIR=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
CURDIR=
DEVIMAGEOPTS=-v $(PROJECT_DIR):/go/src/github.com/fntlnz/netcan -w /go/src/github.com/fntlnz/netcan $(TOOLS_IMAGE_FQDN)
BUILDCONTAINER=docker run --rm $(DEVIMAGEOPTS)
LDFLAGS='-extldflags "-static"' 

ifeq ($(NOCONTAINER), 1)
	BUILDCONTAINER=
	CURDIR=$(PROJECT_DIR)/
endif


build: tools-image
	$(BUILDCONTAINER) sh -c "CGO_ENABLED=0 GOOS=linux go build -ldflags $(LDFLAGS) ."

tools-image:
	docker build -t $(TOOLS_IMAGE_FQDN) -f Dockerfile .
