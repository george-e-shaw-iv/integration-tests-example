FROM golang:1.11-alpine AS src

# Install git
RUN apk add --no-cache git

# Copy Repository
WORKDIR /home/root/go/src/github.com/george-e-shaw-iv/integration-tests-example/
COPY . ./

# Build Go Binary
RUN set -ex; \
    CGO_ENABLED=0 GOOS=linux go build -o ./listd ./cmd/listd;

# Final image, no source code
FROM debian:jessie-slim

# Install Root Ceritifcates
RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        ca-certificates;

WORKDIR /opt/
COPY --from=src /home/root/go/src/github.com/george-e-shaw-iv/integration-tests-example/listd .

# Run Go Binary
CMD /opt/listd