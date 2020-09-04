#!/bin/bash

set -eu

docker run --init \
    --net=host \
    --detach \
    --restart unless-stopped \
    -e BCAST_ADDRESS=192.168.178.255 \
    -p 9132:9132 \
    --name sop112_exporter \
    mrbuk/sop112_exporter:0.3
