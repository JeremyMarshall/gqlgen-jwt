# Image URL to use all building/pushing image targets
IMG ?= jeremymarshall/gqlgen-jwtapi:v0.1.1
PORT ?= 8088
SERVEROPTS ?= dummy

# Build manager binary
api: bin fmt vet
	go build -o bin/api 

server: build
	PORT=${PORT} go run ./main.go ${SERVEROPTS}

gqlgen: graph/resolver.go

graph/resolver.go: graph/*.graphqls gqlgen.yml
	go run github.com/99designs/gqlgen --verbose

build: bin tidy gqlgen fmt vet

tidy:
	go mod tidy

fmt:
	go fmt ./...

vet:
	go vet ./...

bin:
	mkdir $@

test:
	go test -count 1  -coverprofile cover.out ./...
	# go test -count 1 -coverpkg $(PKGS) -coverprofile cover.out ./...

cover: test
	go tool cover -html=cover.out

# Build the docker image
docker-build: test
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}	

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy:
	cd config/api && kustomize edit set image api=${IMG}
	kustomize build config/default | kubectl apply -f -

clean:
	-rm -fr dist bin cover.out coverage.txt cp.out

.PHONY: server build test diff fmt vet tidy cover docker-push docker-build gqlgen
