SHELL := /bin/bash

# ==============================================================================
# testing

# expvarmon -ports=":3310" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ==============================================================================
# run

run:
	go run main.go

# ==============================================================================
# building containers

VERSION := 0.0.1	# this should match BUILD_VERSION in config.yml

all: api
api:
	docker build \
		-f zarf/docker/dockerfile.api \
		-t api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
