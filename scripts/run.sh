#!/bin/bash

set -eu

docker run --init \
    --net=host \
    --detach \
    --restart unless-stopped \
    -p 9132:9132 \
    --name sop112_exporter \
    mrbuk/sop112_exporter:0.6 -broadcast 192.168.178.255
