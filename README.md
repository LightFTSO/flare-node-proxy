# Flare Node Proxy

[![Delegate](https://img.shields.io/badge/Delegate-here-orange)](https://lightft.so/delegate)
[![Github](https://img.shields.io/github/followers/LightFTSO?style=social)](https://github.com/LightFTSO)
[![Twitter](https://img.shields.io/twitter/follow/lightFTSO?style=social)](https://twitter.com/lightFTSO)
[![Telegram](https://img.shields.io/badge/Telegram--blue?logo=telegram)](https://t.me/LightFTSO)
![Security](https://github.com/LightFTSO/flare-node-proxy/workflows/Security/badge.svg)
![Linter](https://github.com/LightFTSO/flare-node-proxy/workflows/Linter/badge.svg)


## Description

This app works as a proxy with public Flare API providers in mind. It includes a middleware that blocks transactions to the PriceSubmitter contract
by default.
It doesn't modify, record, or do anything with the request other than rejecting transactions to the PriceSubmitter contract (**0x1000000000000000000000000000000000000003**).

It uses GoFiber's Proxy to forward requests to the specified endpoint, the requests should get there intact, and come back intact too.
And it also uses [goccy/go-json](https://github.com/goccy/go-json) to make JSON marshalling/unmarshalling faster.

## TODO
* Handle proxy to websockets
* Implement whitelisting of preapproved addresses, with optional expiration dates
* Implement better logging (maybe it's not needed?)


## How to use
### Start the application 


```bash
go run app.go --monitor --port :3000 --endpoint http://localhost:9650
```

### Help
```bash
go run app.go --help
```

### Use local container

```
# Clean packages
make clean-packages

# Generate go.mod & go.sum files
make requirements

# Generate docker image
make build

# Generate docker image with no cache
make build-no-cache

# Run the projec in a local container
make up

# Run local container in background
make up-silent

# Run local container in background with prefork
make up-silent-prefork

# Stop container
make stop

# Start container
make start
```

## Production

```bash
docker build -t flare-node-proxy .
docker run -d -p 3000:3000 flare-node-proxy ./app -prod -monitor -endpoint http://localhost:9650
```
