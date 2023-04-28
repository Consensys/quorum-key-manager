---
title: Stores
description: Description of stores
sidebar_position: 2
---

# Stores

A store is a Quorum Key Manager (QKM) component that interfaces with an underlying secure system storage (such as HashiCorp Vault, Azure Key Vault, or AWS KMS) to perform crypto-operations.

The store manager is a QKM component that uses [store manifest](../HowTo/Use-Manifest-File/Store.md) to create and manage stores. Other QKM components use the store manager to access stores and perform crypto-operations.

QKM defines the following store interfaces:

- [Stores](#stores)
  - [Vault](#vault)
  - [Secret store](#secret-store)
  - [Key store](#key-store)
  - [Ethereum store](#ethereum-store)

## Vault

A vault defines the user credentials required to access secure system storage, such as HashiCorp Vault, Azure Key Vault, or AWS KMS.

## Secret store

A secret store allows you to store and access secret values, but doesn't expose any crypto-operations.

Depending on the implementation, a secret store:

- Can manage multiple versions of a secret.
- Has advanced capabilities to delete and recover secrets.

If you have existing secrets in a secure storage system, you must [index](../HowTo/Index-Resources.md) them in your local QKM database in order to use them. Use the [`/secrets`](https://consensys.github.io/quorum-key-manager/#tag/Secrets) REST API endpoint to interact with a secret store.

## Key store

A key store manages keys and enables crypto-operations such as signing and encryption.

A key store can generate and import keys, but doesn't allow access to the private part of a key.

Depending on the implementation, a key store:

- Has advanced capabilities to delete and recover keys.

You can implement a key store to:

- Delegate crypto-operations to an external dependency.
- Use the underlying secret store to perform crypto-operations locally.

If you have existing keys in a secure storage system, you must [index](../HowTo/Index-Resources.md) them in your local QKM database in order to use them. Use the [`/keys`](https://consensys.github.io/quorum-key-manager/#tag/Keys) REST API endpoint to interact with a key store.

## Ethereum store

An Ethereum store manages Ethereum accounts and performs Ethereum-related crypto-operations (for example, signing transactions).

An Ethereum store can generate and import accounts but does not expose the private key of any account.

You can implement an Ethereum store based on an underlying key store to perform signing, while the account store is responsible for performing Ethereum-specific processing, formatting, and encoding.

If you have existing Ethereum accounts in a secure storage system, you must [index](../HowTo/Index-Resources.md) them in your local QKM database in order to use them. Use the [`/ethereum`](https://consensys.github.io/quorum-key-manager/#tag/Ethereum-Account) REST API endpoint to interact with an Ethereum store.
