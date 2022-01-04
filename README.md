[![Website](https://img.shields.io/website?label=documentation&url=https%3A%2F%2Fdocs.quorum-key-manager.consensys.net%2F)](https://docs.quorum-key-manager.consensys.net/)
[![Website](https://img.shields.io/website?url=https%3A%2F%2Fconsensys.net%2Fquorum%2F)](https://consensys.net/quorum/)

[![CircleCI](https://img.shields.io/circleci/build/gh/ConsenSys/quorum-key-manager?token=7062612dcd5a98913aa1b330ae48b6a527be52eb)](https://circleci.com/gh/ConsenSys/quorum-key-manager)
[![codecov](https://codecov.io/gh/ConsenSys/quorum-key-manager/branch/main/graph/badge.svg)](https://codecov.io/gh/ConsenSys/quorum-key-manager)
[![Go Report Card](https://goreportcard.com/badge/github.com/ConsenSys/quorum-key-manager)](https://goreportcard.com/report/github.com/ConsenSys/quorum-key-manager)

# Quorum Key Manager
Quorum Key Manager (QKM) is a key management service developed under the [BSL 1.1](LICENSE) license and written in Go. 

Quorum Key Manager exposes an HTTP API service to manage your secrets, keys and Ethereum accounts. QKM supports the integration with
*AWS Key Management Service*, *Azure Key Vault* and *HashiCorp Vault*. 

In addition, using the JSON-RPC interface of the QKM, you can connect to your Ethereum nodes to sign your transaction using the Ethereum account stored in your secure key vault.

## Useful links

* [Product page](https://consensys.net/quorum/key-manager/)
* [User documentation](http://docs.quorum-key-manager.consensys.net/)
* [REST API reference documentation](https://consensys.github.io/quorum-key-manager/#stable)
* [GitHub Project](https://github.com/ConsenSys/quorum-key-manager)
* [issues](https://github.com/ConsenSys/quorum-key-manager/issues)
* [Changelog](https://github.com/ConsenSys/quorum-key-manager/blob/main/CHANGELOG.md)
* [HashiCorp Vault plugin](https://github.com/ConsenSys/quorum-hashicorp-vault-plugin)
* [Helm Charts](https://github.com/ConsenSys/quorum-key-manager-helm)
* [Kubernetes deployment example](https://github.com/ConsenSys/quorum-key-manager-kubernetes)

## Run QKM

First, define your Quorum Key Manager environment setup using manifest files.
Examples can be found at [`./deps/config/manifests`](./deps/config/manifests). 
More information about how to set up service can be found in [documentation](http://docs.quorum-key-manager.consensys.net/).

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
 
## License

Orchestrate is licensed under the BSL 1.1. Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [quorum-key-manager@consensys.net](mailto:quorum-key-manager@consensys.net) if you need to purchase a license for a production use-case.  
