# Builder
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base AS builder

WORKDIR /go/src/gitlab.com/ignitionrobotics/web/cloudsim
COPY . /go/src/gitlab.com/ignitionrobotics/web/cloudsim

# Install the dependencies without checking for go code
#RUN dep ensure -vendor-only -v
#COPY vendor

# Build app
RUN go install


# Runner
FROM registry.gitlab.com/ignitionrobotics/web/images/cloudsim-base

WORKDIR /app
COPY --from=builder /go/bin/cloudsim .
COPY . .

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["./cloudsim"]

EXPOSE 8001
