---
title: Add a store
description: Add stores using manifest
sidebar_position: 2
---

# Add a store to Quorum Key Manager

You can define a [store](../../Concepts/Stores.md) in a Quorum Key Manager (QKM) [manifest file](Overview.md).

QKM supports the following store interfaces:

- [Add a store to Quorum Key Manager](#add-a-store-to-quorum-key-manager)
  - [Vault](#vault)
    - [HashiCorp](#hashicorp)
    - [Azure Key Vault](#azure-key-vault)
    - [Amazon Key Management Service](#amazon-key-management-service)
  - [Secret store](#secret-store)
  - [Key store](#key-store)
  - [Ethereum store](#ethereum-store)

:::warning

If you have existing Ethereum accounts, keys, or secrets in a secure storage system, you must [index](../Index-Resources.md) them in your local QKM database in order to use them.

:::

## Vault

Use the following fields to configure one or more [vaults](../../Concepts/Stores.md#vault):

- `kind`: _string_ - vault
- `type`: _string_ - supported vault types are `hashicorp`, `azure`, and `aws`
- `name`: _string_ - identifier of the vault
- `allowed_tenants`: _array_ of _strings_ - (optional) list of allowed tenants for this store when using [resource-based access control](../../Concepts/Authorization.md#resource-based-access-control)
- `specs`: _object_ - [configuration object to connect to an underlying vault](#vault-configuration).

```yaml title="Example vault store manifest file"
# Hashicorp secret store manifest
- kind: Vault
  name: hashicorp-vault
  specs:
    mount_point: secret
    address: http://hashicorp:8200
    token: YOUR_TOKEN
    namespace: user1_space
```

If using one of the following vault services, include the corresponding `spec` fields in your manifest.

### HashiCorp

If using a `HashicorpKeys` or `HashicorpSecrets` store:

- `mount_point`: _string_ - secret engine mounting point
- `address`: _string_ - HashiCorp server URL
- `token_path`: _string_ - path to token file
- `token`: _string_ - authorization token
- `namespace`: _string_ - default namespace to store data in HashiCorp

:::note

- `tokenPath` and `token` are mutually exclusive.
- If using a `Hashicorp` to store keys, you must install the [HashiCorp Vault Plugin](https://github.com/ConsenSys/quorum-hashicorp-vault-plugin).

:::

### Azure Key Vault

If using an `AKVKeys` or `AKVSecrets` store:

- `vault_name`: _string_ - connected Azure Key Vault ID
- `tenant_id`: _string_ - Azure Active Directory tenant ID
- `client_id`: _string_ - user client ID
- `client_secret`: _string_ - user client secret

### Amazon Key Management Service

If using an `AWSKeys` or `AWSSecrets` store:

- `access_id`: _string_ - AWS access ID
- `secret_key`: _string_ - AWS secret key
- `region`: _string_ - AWS region
- `debug`: _boolean_ - indicates whether to enable debugging

## Secret store

Use the following fields to configure one or more [secret stores](../../Concepts/Stores.md#secret-store):

- `kind`: _string_ - Store
- `type`: _string_ - secret
- `name`: _string_ - name of the secret store
- `allowed_tenants`: _array_ of _strings_ - (optional) list of allowed tenants for this store when using [resource-based access control](../../Concepts/Authorization.md#resource-based-access-control)
- `specs`: _object_ - [configuration object to selected injected vault](#vault-configuration).

```yaml title="Example secret store manifest file"
# Hashicorp secret store manifest
- kind: Store
  type: secret
  name: my-secret-store
  specs:
    vault: hashicorp-vault
```

## Key store

Use the following fields to configure one or more [key stores](../../Concepts/Stores.md#key-store):

- `kind`: _string_ - Store
- `type`: _string_ - key
- `name`: _string_ - name of the key store
- `allowed_tenants`: _array_ of _strings_ - (optional) list of allowed tenants for this store when using [resource-based access control](../../Concepts/Authorization.md#resource-based-access-control)
- `specs`: _object_ - [configuration object to selected vault or secret store](#vault-configuration).

```yaml title="Example key store manifest file"
# Hashicorp key store manifest
- kind: Store
  type: key
  name: my-key-store
  specs:
    vault: hashicorp-vault

# Local key store manifest
- kind: Store
  type: local-keys
  name: my-key-store
  specs:
    secret_store: my-secret-store
```

## Ethereum store

Use the following fields to configure one or more [Ethereum stores](../../Concepts/Stores.md#ethereum-store):

- `kind`: _string_ - Store
- `type`: _string_ - Ethereum
- `name`: _string_ - name of the Ethereum store
- `allowed_tenants`: _array_ of _strings_ - (optional) list of allowed tenants for this store when using [resource-based access control](../../Concepts/Authorization.md#resource-based-access-control)
- `specs`: _object_ - [configuration object to selected key store](#vault-configuration).

```yaml title="Example Ethereum store manifest file"
# Ethereum store manifest
- kind: Store
  type: ethereum
  name: my-ethereum-store
  specs:
    key_store: hashicorp-keys
```
