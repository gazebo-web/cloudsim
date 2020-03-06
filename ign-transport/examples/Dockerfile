# This Dockerfile creates an image with the ign-transport examples built and
# ready to be used.
# Once run, you need to join the created docker container and run the examples
# eg. /osrf/ign-transport/example/build/publisher_c or /osrf/ign-transport/example/build/subscriber_c
FROM ubuntu:bionic
# we need ubuntu bionic to compile ign_transport

RUN apt-get update && apt-get install -y sudo apt-utils nano vim tar curl build-essential \
  software-properties-common wget ca-certificates git mercurial \
  netcat inetutils-ping iproute2 \
  && git config --global user.name "web-cloudsim"  \
  && git config --global user.email "web-cloudsim@test.org"

# ----------------------------------------------------------------
# ----------------------------------------------------------------
# Gazebo related

# Download dependencies needed to compile ign_transport dev
RUN apt-get update && apt-get install -y gnupg lsb-release cmake pkg-config cppcheck

# Get Gazebo (and ign_transport) dependencies
RUN  echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-stable $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-stable.list \
  && echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-prerelease $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-prerelease.list \
  && echo "deb http://packages.osrfoundation.org/gazebo/ubuntu-nightly $(lsb_release -cs) main" > /etc/apt/sources.list.d/gazebo-nightly.list \
  && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys D2486D2DD83DB69272AFE98867170598AF249743 \
  && apt-get update && apt-get -y install libignition-transport7-dev

# Clone ign_transport code and build (using branch "c_iface")
RUN mkdir -p /osrf && cd /osrf \
  && mkdir -p /root/.ssh && touch /root/.ssh/config && echo -e "Host bitbucket.org\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config \
  && hg clone https://preyna@bitbucket.org/ignitionrobotics/ign-transport \
  && cd ign-transport \
  && hg up c_iface \
  && cd example && mkdir build && cd build && cmake .. && make -j4

RUN mkdir -p /tmp/logs/

CMD ["/bin/sh", "-c", "--", "mkdir -p /var/log/gzserver/test/ && touch /var/log/gzserver/test/testfile.txt && echo 'test' >> /var/log/gzserver/test/testfile.txt && /osrf/ign-transport/example/build/subscriber_c"]

# ----------------------------------------------------------------
# ----------------------------------------------------------------

# Quick tips:
# remove stopped docker containers: docker rm $(docker ps -a -q)
# remove intermediate and <none> docker images: docker images | grep '<none>' | awk '{ print $3; }' | xargs docker rmi
