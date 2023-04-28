---
title: Index resources
description: How to index resources
sidebar_position: 3
---

# Index resources

If you have existing [Ethereum accounts](../Concepts/Stores.md#ethereum-store), [keys](../Concepts/Stores.md#key-store), or [secrets](../Concepts/Stores.md#secret-store) in a secure storage system, you must index (reference) these resources in your local QKM database using the [`sync` subcommand](../Reference/CLI/CLI-Subcommands.md#sync) in order to use them.

Use `sync ethereum` to index Ethereum accounts, `sync keys` to index keys, and `sync secrets` to index secrets. You can specify options [on the command line](#on-the-command-line) or [as environment variables](#as-environment-variables).

:::warning

When you index your resources, the private keys remain in the underlying secure storage system.

The keys are loaded into your local database to generate metadata linked to the underlying system.

:::

## On the command line

Specify the path to the manifest file in which the store is configured using the `--manifest-path` command line option, and the name of the store using the `--store-name` option. Include [any database options](../Reference/CLI/CLI-Syntax.md#db-database) (any options that begin with `--db-`) that apply to your local database.

:::info

Indexing keys from `hashicorp-keys` on the command line:

```bash
key-manager sync keys --manifest-path="/config/default.yml" --store-name="hashicorp-keys" --db-port=8080
```

:::

## As environment variables

You can index resources using environment variables by specifying:

- The manifest file by using the `MANIFEST_PATH` environment variable
- The identifier of the store to index using the `SYNC_STORE_NAME` environment variable
- Optional database connection settings, such as [`DB_DATABASE`](../Reference/CLI/CLI-Syntax.md#db-database)

:::info

Indexing keys from `hashicorp-keys` as environment variables

```text
    MANIFEST_PATH="/config/default.yml"
    SYNC_STORE_NAME="hashicorp-keys"
```

```bash
key-manager sync keys
```

:::
