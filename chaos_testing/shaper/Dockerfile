FROM ubuntu:18.04

RUN apt-get update \
  && apt-get install -y \
      iproute2 \
      ruby2.5 \
  && apt-get clean

COPY shaper.yml ./
COPY shaper.rb ./
ENTRYPOINT ["./shaper.rb"]
