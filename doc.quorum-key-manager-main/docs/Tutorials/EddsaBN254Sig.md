---
title: Sign with EDDSA and BN254
description: Sign with eddsa and babyjubjub
sidebar_position: 2
---

# Sign a transaction with EDDSA and Baby Jubjub

This tutorial walks you through signing an Ethereum transaction with Quorum Key Manager (QKM) using the EDDSA signing algorithm and Baby Jubjub elliptic curve (also referred to as the BN254 twisted Edwards curve).

## Prerequisites

- [Quorum Key Manager installed](../Get-Started/Build-From-Source.md)
- [`curl` command line](https://curl.se/download.html)
- [HashiCorp Vault](https://github.com/hashicorp/vault) service running
- [HashiCorp Vault - Orchestrate Plugin](https://github.com/ConsenSys/orchestrate-hashicorp-vault-plugin) loaded in HashiCorp Vault service

## Steps

1.  In the QKM [manifest file](../HowTo/Use-Manifest-File/Overview.md), specify an [Ethereum store](../HowTo/Use-Manifest-File/Store.md#ethereum-store) to allocate your Ethereum wallets, and the [RPC node](../HowTo/Use-Manifest-File/Node.md) to proxy your calls using QKM.

    ```yaml title="Example manifest file"
    - kind: HashicorpKeys
      name: hashicorp-keys
      specs:
        mountPoint: "{ENGINE_MOUNT_POINT}"
        address: "{HASHICORP_VAULT_URL}"
        tokenPath: "{VAULT_TOKEN_PATH}"
        namespace: "{KEYS_NAMESPACE}"

    - kind: Node
      name: besu-node
      specs:
        rpc:
          addr: http://besu-node:8545
    ```

2.  Start QKM with the manifest file by using the [`--manifest-path`](../Reference/CLI/CLI-Syntax.md#manifest-path) option:

    ```bash
    key-manager run --manifest-path=<PATH-TO-MANIFEST-FILE>
    ```

3.  Create an Ethereum account using EDDSA and Baby Jubjub:

    <!--tabs-->

    # curl HTTP request

    ```bash
    curl --request POST 'http://localhost:8080/stores/hashicorp-keys/keys/bn254-key' --header 'Content-Type: application/json' --data-raw '{"curve": "babyjubjub", "signingAlgorithm": "eddsa"}'
    ```

    # JSON result

    ```json
    {
      "id": "bn254-key",
      "publicKey": "Cjix/fS3WdqKGKabagBNYwcClan5aImoFpnjSF0cqJs=",
      "curve": "babyjubjub",
      "signingAlgorithm": "eddsa",
      "disabled": false,
      "createdAt": "2021-09-09T11:18:51.5877561Z",
      "updatedAt": "2021-09-09T11:18:51.5877561Z"
    }
    ```

    <!--/tabs-->

4.  Sign a payload using the created key pair:

    <!--tabs-->

    # Generate base64 message to sign

    ```bash
    echo -n "my signed message" | base64
    ```

# Base64 encoding result

    ```bash
    bXkgc2lnbmVkIG1lc3NhZ2U=
    ```

# curl HTTP request

    ```bash
    curl --request POST 'http://localhost:8080/stores/hashicorp-keys/keys/bn254-key/sign' --header 'Content-Type: application/json' --data-raw '{"data": "bXkgc2lnbmVkIG1lc3NhZ2U="}'
    ```

# JSON result

    ```json
    tjThYhKSFSKKvsR8Pji6EJ+FYAcf8TNUdAQnM7MSwZEEaPvFhpr1SuGpX5uOcYUrb3pBA8cLk8xcbKtvZ56qWA==
    ```

    <!--/tabs-->

5.  Verify your message:

    <!--tabs-->

    # curl HTTP request

    ```bash
    curl --request POST 'http://localhost:8080/stores/hashicorp-keys/keys/verify-signature' --header 'Content-Type: application/json' --data-raw '{"curve": "babyjubjub", "signingAlgorithm": "eddsa", "data": "bXkgc2lnbmVkIG1lc3NhZ2U=", "publicKey": "yhUiySkg/cKbiN8soKZ5YO0GXHqzx8iycnABzYMPE5A=", "signature": "tjThYhKSFSKKvsR8Pji6EJ+FYAcf8TNUdAQnM7MSwZEEaPvFhpr1SuGpX5uOcYUrb3pBA8cLk8xcbKtvZ56qWA=="}'
    ```

    <!--/tabs-->
