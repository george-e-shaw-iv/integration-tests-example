FROM golang:1.12-alpine AS src

# Install git
RUN set -ex; \
    apk update; \
    apk add --no-cache git

# Copy Repository
WORKDIR /go/src/github.com/george-e-shaw-iv/integration-tests-example/
COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -o ./listd ./cmd/listd;

# Final image, no source code
FROM alpine:latest

# Install Root Ceritifcates
RUN set -ex; \
    apk update; \
    apk add --no-cache \
     ca-certificates

WORKDIR /opt/
COPY --from=src /go/src/github.com/george-e-shaw-iv/integration-tests-example/listd .

# Run Go Binary
CMD /opt/listd