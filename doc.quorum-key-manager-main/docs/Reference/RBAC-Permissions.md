---
title: RBAC permissions
description: Role-based access control permissions
sidebar_position: 3
---

# Role-based access control permissions

The following tables list the permissions for [role-based access control](../Concepts/Authorization.md#role-based-access-control). Each permission has a list of allowed [REST endpoints](Rest.md).

## Ethereum accounts

| Name | Description | Allowed endpoints |
| --: | --: | --: |
| `read:ethereum` | Allows reading operations over Ethereum accounts | Get, list, get deleted, list deleted |
| `write:ethereum` | Allows creating Ethereum accounts | Create, import, update |
| `delete:ethereum` | Allows soft-deleting Ethereum accounts | Delete, restore |
| `destroy:ethereum` | Allows permanently deleting Ethereum accounts | Delete, restore, destroy |
| `sign:ethereum` | Allows signing and verifying signatures | _All sign endpoints_, EC recover |
| `encrypt:ethereum` | Allows encryption and decryption | Encrypt, decrypt |

## Keys

| Name | Description | Allowed endpoints |
| --: | --: | --: |
| `read:key` | Allows reading operations over keys | Get, list, get deleted, list deleted |
| `write:key` | Allows creating keys | Create, import, update |
| `delete:key` | Allows soft-deleting keys | Delete, restore |
| `destroy:key` | Allows permanently deleting keys | Delete, restore, destroy |
| `sign:key` | Allows signing and verifying signatures | Sign |
| `encrypt:key` | Allows encryption and decryption | Encrypt, decrypt |
|  |

## Secrets

| Name | Description | Allowed endpoints |
| --: | --: | --: |
| `read:secret` | Allows reading operations over secrets | Get, list, get deleted, list deleted |
| `write:secret` | Allows creating secrets | Set, update |
| `delete:secret` | Allows soft-deleting secrets | Delete, restore |
| `destroy:secret` | Allows permanently deleting secrets | Delete, restore, destroy |

## Alias

| Name           | Description                            | Allowed endpoints |
| :------------- | :------------------------------------- | :---------------- |
| `read:alias`   | Allows reading aliases over registries | Get, list         |
| `write:alias`  | Allows creating aliases                | Create, update    |
| `delete:alias` | Allows deleting aliases                | Delete            |

## Nodes

|          Name |                            Description |
| ------------: | -------------------------------------: |
| `proxy:nodes` | Allows you to proxy traffic into nodes |
