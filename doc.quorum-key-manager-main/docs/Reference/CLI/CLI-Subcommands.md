---
title: Subcommands
description: Quorum Key Manager command line interface subcommands
sidebar_position: 2
---

# Quorum Key Manager subcommands

This reference describes the syntax of the Quorum Key Manager (QKM) command line interface (CLI) subcommands.

You can specify QKM subcommands on the command line:

```bash
key-manager [COMMAND] [SUBCOMMAND] [OPTIONS]
```

## `sync`

Locally [indexes resources](../../HowTo/Index-Resources.md) from a store.

### `ethereum`

<!--tabs-->

# Syntax

```bash
key-manager sync ethereum --manifest-path=<PATH> --store-name=<STRING> [--db-*]
```

# Example

```bash
key-manager sync ethereum --manifest-path="/config/default.yml" --store-name="eth-accounts" --db-database="postgres" --db-host=127.0.0.1 --db-port=6174
```

<!--/tabs-->

Indexes (adds references to) Ethereum accounts from the specified [Ethereum store](../../Concepts/Stores.md#ethereum-store) configured in the specified [manifest file](../../HowTo/Use-Manifest-File/Overview.md) into your local QKM database.

You must specify:

- The path of the manifest file using the `--manifest-path` option or the `MANIFEST_PATH` environment variable.
- The name of the store from which to index the resources using the `--store-name` option or the `STORE_NAME` environment variable.

You can include any [database options or environment variables](CLI-Syntax.md#db-database) (any options that begin with `--db-`).

### `keys`

<!--tabs-->

# Syntax

```bash
key-manager sync keys --manifest-path=<PATH> --store-name=<STRING> [--db-*]
```

# Example

```bash
key-manager sync keys --manifest-path="/config/default.yml" --store-name="hashicorp-keys" --db-database="postgres" --db-host=127.0.0.1 --db-port=6174
```

<!--/tabs-->

Indexes (adds references to) keys from the specified [key store](../../Concepts/Stores.md#key-store) configured in the specified [manifest file](../..//HowTo/Use-Manifest-File/Overview.md) into your local QKM database.

You must specify:

- The path of the manifest file using the `--manifest-path` option or the `MANIFEST_PATH` environment variable.
- The name of the store from which to index the resources using the `--store-name` option or the `STORE_NAME` environment variable.

You can include any [database options or environment variables](CLI-Syntax.md#db-database) (any options that begin with `--db-`).

### `secrets`

<!--tabs-->

# Syntax

```bash
key-manager sync secrets --manifest-path=<PATH> --store-name=<STRING> [--db-*]
```

# Example

```bash
key-manager sync secrets --manifest-path="/config/default.yml" --store-name="hashicorp-secrets" --db-database="postgres" --db-host=127.0.0.1 --db-port=6174
```

<!--/tabs-->

Indexes (adds references to) secrets from the specified [secret store](../../Concepts/Stores.md#ethereum-store) configured in the specified [manifest file](../../HowTo/Use-Manifest-File/Overview.md) into your local QKM database.

You must specify:

- The path of the manifest file using the `--manifest-path` option or the `MANIFEST_PATH` environment variable.
- The name of the store from which to index the resources using the `--store-name` option or the `STORE_NAME` environment variable.

You can include any [database options or environment variables](CLI-Syntax.md#db-database) (any options that begin with `--db-`).
