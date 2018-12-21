FROM golang:1.11

# Copy Repository
WORKDIR /integration-tests-example/
COPY . ./

# Install Root Ceritifcates
RUN set -ex; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
        ca-certificates;

# Build Go Binary
RUN set -ex; \
    GO111MODULE=on go get ./...; \
    GO111MODULE=on go build ./...; \
    CGO_ENABLED=0 GOOS=linux go build -o /opt/listd ./cmd/listd;

# Run Go Binary
CMD /opt/listd