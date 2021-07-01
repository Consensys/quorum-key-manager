# Quorum Key Manager
Quorum Key Manager(QKM) is a key management service developed under the BSL 1.1 license and written in Go. 

Quorum Key Manager exposes a HTTP API service to manage your secrets, keys and ethereum accounts. QKM supports the integration with
*Amazon Key Management Service*, *Azure Key Vault* and *Hashicorp Vault*. 

In addition, using QKM, you can connect to your ethereum nodes to sign your transaction using the ethereum account stored in your secure key vault.

## Run QKM

Available docker images can be found in `docker.consensys.net/priv/quorum-key-manager`.

To run the Quorum Key Manager service using docker you can execute the following command:

```
docker run -it \
--name quorum-key-manager \
--mount  type=bind,source="$(pwd)"/deps/config,target=/manifests \
docker.consensys.net/priv/quorum-key-manager:stable run --manifest-path=/manifests
```

You can find more information about the expected content of the `/manifest` folder in the project [documentation](#documentation) 

## Build binaries

To build binary locally requires Go (version 1.15 or later) and C compiler. 

After installing project vendors (ie `go mod vendor`) you can run following command to compile the binary

```
make gobuild
```

Binary will be located in `./build/bin/key-manager`

## Documentation

Quorum Key Manager documentation website [https://docs.quorum-key-manager.consensys.net/](https://docs.quorum-key-manager.consensys.net/) 

 
## License

Orchestrate is licensed under the BSL 1.1.

Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [orchestrate@consensys.net](mailto:orchestrate@consensys.net) if you need to purchase a license for a
production use-case.

## AWS specific annotations

When getting one Key, KM may return the following annotations :

- "aws-KeyID" : the Key inner ID
- "aws-KeyStoreID" : the Key's Keystore ID
- "aws-ClusterHSMID" : the Key's Cluster HSM ID, when backed by an actual HSM
- "aws-AccountID" : the current user (requester) account ID
- "aws-ARN" : the key Amazon Ressource Name