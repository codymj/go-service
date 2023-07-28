SHELL := /bin/bash

# ==============================================================================
# testing

# docker run -p 3310:3310 -p 3300:3300 <imgId>
# docker image prune -a
# expvarmon -ports=":3310" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
# hey -m GET -c 100 -n 10000 http://localhost:3300/v1/test

# generate public/private keys
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem

# test query
# curl --user "admin@example.com:gophers" http://localhost:3300/v1/users/token
# curl -H "Authorization: Bearer ${TOKEN} http://localhost:3300/v1/users/1/2

# test auth
# do "make admin" and "export TOKEN=<token>"
# curl -w "\n" -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3300/v1/testauth

# dblab --host 127.0.0.1 --user postgres --db postgres --pass <password> --schema public --ssl disable --port 5432 --driver postgres --limit 50

# ==============================================================================

# run service
run:
	go run app/services/api/main.go

# creates a x509 public/private key pair for auth tokens
genkeys:
	go run app/tools/admin/main.go genkeys

# generates jwt token
gentoken:
	go run app/tools/admin/main.go gentoken

# migrates database
migrate:
	go run app/tools/admin/main.go migrate

# seeds database with test data
seed:
	go run app/tools/admin/main.go seed

# run all tests
test:
	go test ./... -count=1

# update service
update:
	go get -t -u ./... && go mod vendor

# ==============================================================================
# building containers

# this should match BUILD_VERSION in config.yml
VERSION := 0.0.1

all: api
api:
	docker build \
		-f zarf/docker/dockerfile.api \
		-t api-amd64:$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
