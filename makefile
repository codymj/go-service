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

# ==============================================================================
# running from within k8s/kind

KIND_CLUSTER := cluster0

kind-up:
	kind create cluster \
		--image kindest/node:v1.26.0@sha256:3264cbae4b80c241743d12644b2506fff13dce07fcadf29079c1d06a47b399dd \
		--name $(KIND_CLUSTER)
		--config zarf/k8s/kind/config.yaml

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image service-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
