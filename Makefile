
.PHONY: build deploy clean

GOARCH ?= $(shell go env GOARCH)

GOBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=$(GOARCH) go build -mod=vendor
build:
	$(GOBUILD) -o bin/mycnid cmd/mycnid/main.go
	$(GOBUILD) -o bin/mycni cmd/mycni/main.go

VERSION=v1.0
docker-build: build
	docker build -t mycni:$(VERSION) .

deploy:
	kubectl apply -f deploy/mycni.yaml

clean:
	rm -rf bin
	go mod tidy
	go mod vendor

kind-cluster:
	kind create cluster --config deploy/kind.yaml

kind-image-load: docker-build
	kind load docker-image mycni:$(VERSION) 
