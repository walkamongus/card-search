# Build the app binary
FROM golang:1.16 as builder

WORKDIR /workspace

# Copy in module manifests and download required modules
COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

# Copy in the source files
COPY main.go main.go
COPY handlers.index.go handlers.index.go
COPY routes.go routes.go
COPY internal/ internal/
COPY templates/ templates/

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o app .

# Package up binary in a minimal distroless container
# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/app .
COPY --from=builder /workspace/templates/ templates/
USER nonroot:nonroot
ENTRYPOINT ["/app"]
