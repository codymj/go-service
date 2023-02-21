SHELL := /bin/bash

run:
	go run main.go

# ==============================================================================
# building containers

VERSION := 0.0.1

all: service
service:
	docker build \
		-f zarf/docker/dockerfile \
		-t service-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
