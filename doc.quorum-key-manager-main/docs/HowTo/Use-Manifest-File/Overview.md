---
title: Overview
description: Using a manifest file
sidebar_position: 1
---

# Using the Quorum Key Manager manifest file

Use a YAML manifest file to configure the Quorum Key Manager (QKM) runtime components. You can configure:

- [Stores](Store.md) - A store interfaces with an underlying secure system storage (such as HashiCorp Vault, Azure Key Vault, or AWS KMS) to perform crypto-operations.
- [Nodes](Node.md) - A node interfaces with underlying node endpoints (such as RPC nodes and Tessera nodes).
- [Roles](Role.md) - A role is a named set of permissions assigned to a user.

You can define multiple manifests in one manifest file, each separated by a dash (`-`).

```yaml title="Example Quorum Key Manager manifest file"
# Hashicorp secret store manifest
- kind: Vault
  type: hashicorp
  name: hashicorp-vault
  specs:
    mount_point: secret
    address: http://hashicorp:8200
    token_path: path/to/token_file
    token: YOUR_TOKEN
    namespace: user1_space

- kind: Store
  type: secret
  name: hashicorp-secrets
  specs:
    vault: hashicorp-vault

# GoQuorum node manifest
- kind: Node
  name: goquorum-node
  specs:
    rpc:
      addr: http://goquorum1:8545
    tessera:
      addr: http://tessera1:9080
```

Specify the path to the manifest file or to a directory with several manifest files using the [`--manifest-path`](../../Reference/CLI/CLI-Syntax.md#manifest-path) command line option on QKM startup. You can alternatively use the `MANIFEST_PATH` environment variable.

```bash title="Starting Quorum Key Manager with a manifest file"
key-manager run --manifest-path=/config/manifest.yml
```
