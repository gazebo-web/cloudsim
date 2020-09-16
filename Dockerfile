# Builder
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base:1.1.0 AS builder

WORKDIR /go/src/gitlab.com/ignitionrobotics/web/cloudsim
COPY . /go/src/gitlab.com/ignitionrobotics/web/cloudsim

# Build app
RUN go install

# Runner
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base:1.1.0

WORKDIR /app
COPY --from=builder /go/bin/cloudsim .
COPY . .

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["./cloudsim"]

EXPOSE 8001
