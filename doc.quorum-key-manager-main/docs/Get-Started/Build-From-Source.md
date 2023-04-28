---
description: Building Quorum Key Manager from source code
sidebar_position: 2
---

# Build from source

## Prerequisites

- [Go version 1.15 or later](https://golang.org/doc/install)
- C compiler such as [GCC](https://gcc.gnu.org/)
- [Postgres](https://www.postgresql.org/)

## Installation on Linux / Unix / macOS

Clone the `Consensys/quorum-key-manager` repository:

```bash
git clone https://github.com/Consensys/quorum-key-manager.git
```

Go to the `quorum-key-manager` directory:

```bash
cd quorum-key-manager
```

Install the project vendors:

```bash
go mod download
```

Compile the binary:

```bash
go build -o ./build/bin/key-manager
```

Display help information and confirm installation:

```bash
./build/bin/key-manager run --help
```

You can optionally create an alias:

```bash
alias key-manager="<PATH-TO-QUORUM-KEY-MANAGER>/build/bin/key-manager"
```

### Start Quorum Key Manager

Specify environment variables to connect to Postgres:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_DATABASE=qkm
```

Run all database migrations:

```bash
key-manager migrate up
```

:::warning

If you have existing [Ethereum accounts](../Concepts/Stores.md#ethereum-store), [keys](../Concepts/Stores.md#key-store), or [secrets](../Concepts/Stores.md#secret-store) in a secure storage system (such as HashiCorp KV, AWS KMS, or Azure KV), in order to use those resources you must [index](../HowTo/Index-Resources.md) them after running all database migrations and before or after starting Quorum Key Manager.

:::

Start Quorum Key Manager specifying the path to a [manifest file or folder](../HowTo/Use-Manifest-File/Overview.md) and any other [options](../Reference/CLI/CLI-Syntax.md):

```bash
key-manager run --manifest-path=<PATH> [OPTIONS]
```
