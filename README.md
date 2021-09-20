[![Website](https://img.shields.io/website?label=documentation&url=https%3A%2F%2Fdocs.quorum-key-manager.consensys.net%2F)](https://docs.quorum-key-manager.consensys.net/)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Fconsensys.net%2Fquorum%2F)](https://consensys.net/quorum/)

[![CircleCI](https://img.shields.io/circleci/build/gh/ConsenSys/quorum-key-manager?token=7062612dcd5a98913aa1b330ae48b6a527be52eb)](https://circleci.com/gh/ConsenSys/quorum-key-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/ConsenSys/quorum-key-manager)](https://goreportcard.com/report/github.com/ConsenSys/quorum-key-manager)

# Quorum Key Manager
Quorum Key Manager (QKM) is a key management service developed under the [BSL 1.1](LICENSE) license and written in Go. 

Quorum Key Manager exposes an HTTP API service to manage your secrets, keys and Ethereum accounts. QKM supports the integration with
*AWS Key Management Service*, *Azure Key Vault* and *HashiCorp Vault*. 

In addition, using QKM, you can connect to your Ethereum nodes to sign your transaction using the Ethereum account stored in your secure key vault.

## Run QKM

First, define your Quorum Key Manager environment setup using manifest files.
Examples can be found at [`./deps/config/manifests`](./deps/config/manifests). 
More information about how to setup service can be found in [documentation](#documentation).

Once you manifests are written, specify where `quorum-key-manager` can find them on your filesystem
with the following environment variable: 

```bash
export HOST_MANIFEST_PATH={your_manifests_folder}
```
 
Now launch Quorum Key Manager service using docker-compose with the following command:

```bash
docker-compose up
```

## Build from source

### Prerequisites

To build binary locally requires Go (version 1.15 or later) and C compiler. 

### Build

After downloading dependencies (ie `go mod download`) you can run following command to compile the binary

```bash
go build -o ./build/bin/key-manager
```

Binary will be located in `./build/bin/key-manager`

## Documentation

Quorum Key Manager documentation website [https://docs.quorum-key-manager.consensys.net/](https://docs.quorum-key-manager.consensys.net/) 
 
## License

Orchestrate is licensed under the BSL 1.1. Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [quorum-key-manager@consensys.net](mailto:quorum-key-manager@consensys.net) if you need to purchase a license for a production use-case.  
