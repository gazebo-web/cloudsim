#!/bin/bash
set -e
set -x

if [ ! "$IGN_TRANSPORT_IP_INTERFACE" = "" ]; then
    # extract IP
    ip=$(ifconfig $IGN_TRANSPORT_IP_INTERFACE | sed -En 's/127.0.0.1//;s/.*inet (addr:)?(([0-9]*\.){3}[0-9]*).*/\2/p')
    export IGN_IP="$ip"
fi

exec "$@"
