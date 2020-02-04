FROM ubuntu:bionic
RUN apt-get update \
      && apt-get install -y \
      sudo \
      curl \
      python3-pip \
      wget \
      vim \
      build-essential \
      apt-transport-https \
      ca-certificates \
      gnupg-agent \
      groff \
      software-properties-common

# Install docker
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
RUN apt-get update
RUN apt-get install -y docker-ce docker-ce-cli containerd.io

# Install AWS CLI
RUN pip3 install --upgrade pip
RUN pip3 install awscli --upgrade

# Configure some default values. This is not the recommend method to set
# secret information. Looking at the docker history will reveal the secret keys
# We are using this method for demonstration purposes only.
ARG access_key_id
ARG secret_access_key
RUN aws configure set aws_access_key_id $access_key_id
RUN aws configure set aws_secret_access_key $secret_access_key
RUN aws configure set default.region us-east-1
RUN aws configure set region us-east-1

