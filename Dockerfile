FROM ignitionrobotics/web-cloudsim-base

RUN mkdir -p /go/src/gitlab.com/ignitionrobotics/web/cloudsim
COPY . /go/src/gitlab.com/ignitionrobotics/web/cloudsim
WORKDIR /go/src/gitlab.com/ignitionrobotics/web/cloudsim

# Install the dependencies without checking for go code
RUN dep ensure -vendor-only

# Build app
RUN go install

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["/go/bin/cloudsim"]

EXPOSE 8001
