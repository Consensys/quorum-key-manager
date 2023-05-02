---
title: Connect to the JSON-RPC node proxy
description: JSON-RPC node proxy tutorial
sidebar_position: 1
---

# Connect to the JSON-RPC node proxy

This tutorial walks you through connecting to the JSON-RPC node proxy and signing an Ethereum transaction using Quorum Key Manager (QKM) as remote and secure storage for your wallets.

## Prerequisites

- [Quorum Key Manager installed](../Get-Started/Build-From-Source.md)
- [`curl` command line](https://curl.se/download.html)
- [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/) configured

## Steps

1.  In the QKM [manifest file](../HowTo/Use-Manifest-File/Overview.md), specify an [Ethereum store](../HowTo/Use-Manifest-File/Store.md#ethereum-store) to allocate your Ethereum wallets, and the [RPC node](../HowTo/Use-Manifest-File/Node.md) to proxy your calls using QKM.

    ```yaml title="Example manifest file"
    - kind: Vault
      type: azure
      name: akv-europe
      specs:
        keystore: AzureKeys
        specs:
          vaultName: <AZURE-VAULT-ID>
          tenantID: <TENANT-ID>
          clientID: <CLIENT-ID>
          clientSecret: <SECRET>

    - kind: Store
      type: key
      name: akv-keys
      specs:
        vault: akv-europe

    - kind: Store
      type: ethereum
      name: eth-accounts
      specs:
        key_store: akv-keys

    - kind: Node
      name: quorum-node
      specs:
        rpc:
          addr: http://quorum1:8545
        tessera:
          addr: http://tessera1:9080
    ```

2.  Start QKM with the manifest file by using the [`--manifest-path`](../Reference/CLI/CLI-Syntax.md#manifest-path) option:

    ```bash
    key-manager run --manifest-path=<PATH-TO-MANIFEST-FILE>
    ```

3.  Create an Ethereum account:

    <!--tabs-->

    # curl HTTP request

    ```bash
    curl -X POST 'http://localhost:8080/stores/eth-accounts/ethereum'
    ```

    # JSON result

    ```json
    {
      "publicKey": "0x045c36d8acc9b00a33221cea6caa39a826e396c0be9df00d224c7aa077b4b58a18e6fdf79a4e9724f9f61a8cdac691c3fea30309be0f46035e299051e4c95a62b3",
      "compressedPublicKey": "0x035c36d8acc9b00a33221cea6caa39a826e396c0be9df00d224c7aa077b4b58a18",
      "createdAt": "2021-07-02T07:33:26.24350701Z",
      "updatedAt": "2021-07-02T07:33:26.24350701Z",
      "expireAt": "0001-01-01T00:00:00Z",
      "keyId": "qkm--95BdaiyQ8OEyX8a",
      "address": "0xd8c88f28748367a11d3c6fc010eef7b670ac016f",
      "disabled": false
    }
    ```

    <!--/tabs-->

4.  Sign a transaction using the Ethereum account and the RPC node proxy:

    <!--tabs-->

    # curl HTTP request

    ```bash
    curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from": "0xd8c88f28748367a11d3c6fc010eef7b670ac016f","to": "0xd46e8dd67c5d32be8058bb8eb970870f07244567", "data":"0xafed"}], "id":1}' http://localhost:8080/nodes/quorum-node
    ```

    # JSON result

    ```json
    {
      "jsonrpc": "2.0",
      "result": "0x8c961ba2c3f51f9088e1a12a81bb1ad9c551ccfad75615f39e4fc95c3bb7086b",
      "error": null,
      "id": 1
    }
    ```

    <!--/tabs-->

5.  Fetch the transaction receipt using the RPC node proxy:

    <!--tabs-->

    # curl HTTP request

    ```bash
    curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x8c961ba2c3f51f9088e1a12a81bb1ad9c551ccfad75615f39e4fc95c3bb7086b"],"id":1}' http://localhost:8080/nodes/quorum-node
    ```

    # JSON result

        ```json
        {"jsonrpc":"2.0","result":{"blockHash":"0x593a660cbd41df2bb58e56bdb70265c8d2738e5d8c9f01bd47e10eec89ebe052","blockNumber":"0x9b","contractAddress":null,"cumulativeGasUsed":"0x5290","from":"0xf772512f135c92a94a0fece58222c982bec0b837","gasUsed":"0x5290","logs":[],"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","status":"0x1","to":"0xd46e8dd67c5d32be8058bb8eb970870f07244567","transactionHash":"0xe684bfe231e9b4d2f3b309532a495bcf8a9acc369940b0bd464678987f1276a3","transactionIndex":"0x0"},"error":null,"id":1}
        ```

    <!--/tabs-->
