# Build the api binary
FROM golang:1.12.5 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY graph/ graph/
COPY rbac/ rbac/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o api main.go

# Use distroless as minimal base image to package the api binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot 
WORKDIR /
COPY --from=builder /workspace/api .
COPY all.yaml .
USER nonroot:nonroot

ENTRYPOINT ["/api"]
