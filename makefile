SHELL := /bin/bash

run:
	go run main.go

# ==============================================================================
# building containers

VERSION := 0.0.1

all: api
api:
	docker build \
		-f zarf/docker/dockerfile.api \
		-t api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
