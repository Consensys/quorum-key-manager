---
title: Add a node
description: Add nodes using manifest
sidebar_position: 3
---

# Add a node to Quorum Key Manager

You can define a [node](../../Concepts/Nodes.md) in a Quorum Key Manager (QKM) [manifest file](Overview.md). For each defined node, QKM exposes HTTP and WebSocket endpoints to interact with the node, using [Ethereum stores](../../Concepts/Stores.md#ethereum-store) as remote and secure key stores.

Use the following fields to configure one or more nodes:

- `kind`: _string_ - the string `Node`
- `name`: _string_ - name of the node
- `allowed_tenants`: _array_ of _strings_ - (optional) list of allowed tenants for this node when using [resource-based access control](../../Concepts/Authorization.md#resource-based-access-control)
- `specs`: _object_ - configuration object to connect to various endpoints, with the following fields for each endpoint:
  - `rpc` or `tessera`: (field name is the name of the endpoint)
    - `addr`: _string_ - address of the endpoint
- `tags`: _map_ of _strings_ to _strings_ - (optional) user set information about the node

:::info

```yaml title="Example node manifest file"
# GoQuorum node manifest
- kind: Node
  name: goquorum-node
  specs:
    rpc:
      addr: http://goquorum1:8545
    tessera:
      addr: http://tessera1:9080

# Besu node manifest
- kind: Node
  name: besu-node
  specs:
    rpc:
      addr: http://validator1:8545
```

:::
