FROM ubuntu:bionic

# Install dependencies
RUN apt-get update && apt-get install -y sudo apt-utils nano vim tar curl build-essential \
  iproute2 inetutils-ping net-tools \
  software-properties-common wget ca-certificates git mercurial

# Config git

RUN git config --global user.name "web-cloudsim"
RUN git config --global user.email "web-cloudsim@test.org"

# Install Gazebo

# Download dependencies needed to compile ign_transport dev
RUN apt-get update && apt-get install -y gnupg lsb-release cmake pkg-config cppcheck

# Get Gazebo (and ign_transport) dependencies
RUN  echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-stable $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-stable.list \
  && echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-nightly $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-nightly.list \
  && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys D2486D2DD83DB69272AFE98867170598AF249743 \
  && apt-get update && apt-get -y install libignition-transport7-dev

# Download and install Go 1.10.3
# More details here: https://github.com/docker-library/golang/blob/master/1.11/stretch/Dockerfile
RUN cd ~ && curl -O https://dl.google.com/go/go1.10.3.linux-amd64.tar.gz  \
  && tar -C /usr/local -xzf go1.10.3.linux-amd64.tar.gz
ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

# Install go dep (v.0.4.1)
ENV DEP_RELEASE_TAG=v0.4.1
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Download and install kubectl (v1.13.1)
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/v1.13.1/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin/kubectl

########################################################################################################################

# Create Kube folder
RUN mkdir -p /root/.kube

RUN mkdir -p /go/src/gitlab.com/ignitionrobotics/web/cloudsim
COPY . /go/src/gitlab.com/ignitionrobotics/web/cloudsim
WORKDIR /go/src/gitlab.com/ignitionrobotics/web/cloudsim

# Install the dependencies without checking for go code
RUN dep ensure -vendor-only

# Build app
RUN go install

# Copy kube config file to .kube folder
COPY kube_config /root/.kube/config

ENTRYPOINT [ "./docker-entrypoint.sh" ]
CMD ["/go/bin/cloudsim"]

EXPOSE 8001
