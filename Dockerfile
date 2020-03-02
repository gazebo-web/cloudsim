FROM ignitionrobotics/web-cloudsim-base

RUN mkdir -p /go/src/bitbucket.org/ignitionrobotics/web-cloudsim
COPY . /go/src/bitbucket.org/ignitionrobotics/web-cloudsim
WORKDIR /go/src/bitbucket.org/ignitionrobotics/web-cloudsim

# Install the dependencies without checking for go code
RUN dep ensure -vendor-only

# Build app
RUN go install

# Copy kube config file to .kube folder
COPY kube_config /root/.kube/config

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["/go/bin/web-cloudsim"]

EXPOSE 8001
