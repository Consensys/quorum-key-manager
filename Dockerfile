ARG VERSION=nonroot

# Build the manager binary
FROM golang:1.16 as builder

RUN apt-get update && \
	apt-get install --no-install-recommends -y \
	upx-ucl

WORKDIR /app
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o /bin/main main.go
RUN upx /bin/main

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:$VERSION
WORKDIR /
COPY --from=builder /bin/main .
COPY ./deps/migrations /migrations
COPY LICENSE .
USER 65532:65532

ENTRYPOINT ["/main"]
