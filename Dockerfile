# Builder
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base:1.1.0 AS builder

# Copy the source code
WORKDIR /go/src/gitlab.com/ignitionrobotics/web/cloudsim

# Get dependencies
# This step is done explicitly to allow caching this layer
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build app
RUN go install -o cloudsim ./cmd/subt

# Runner
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base:1.1.0

WORKDIR /app

COPY --from=builder /go/bin/cloudsim .
COPY . .

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["./cloudsim"]

EXPOSE 8001
