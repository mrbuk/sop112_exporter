# sop112_exporter
Prometheus exporter for WiFi enabled power socket [SOP112](https://de.elv.com/elv-spar-set-mit-2x-schaltsteckdose-sop112-und-1x-powerline-adapter-118025).

For now exports `powersocket_consumption_watts` metrics for each `device` detected

```
powersocket_consumption_watts{device="SWP1040003000001"} 0
powersocket_consumption_watts{device="SWP1040003000002"} 142.74
```

Devices are detected using a specific message to the broadcast IP. The broadcast IP is not auto-discovered and needs to be provided.

## Usage

Run via the CLI

```
$ ./build/sop112_exporter
Usage of ./build/sop112_exporter:
  -broadcast string
        broadcast address to be used for device detected (required)
  -listen string
        address to expose metrics endpoint (default ":9132")
```

or use as a docker container (**needs host mode for the broadcast message sent during detection phase**)

```
docker run \
    --net=host \
    --detach \
    --restart unless-stopped \
    -e BCAST_ADDRESS=192.168.178.255 \
    -p 9132:9132 \
    --name sop112_exporter \
    mrbuk/sop112_exporter:0.1
```

or use `./scripts/run_docker.sh`
