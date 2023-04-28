---
title: User authorization
description: Authorization concept page
sidebar_position: 5
---

# User authorization

After Quorum Key Manager (QKM) [authenticates](Authentication.md) an incoming request, it submits the request to the targeted service which performs authorization checks based on request context before performing service operations.

The authorization process restricts system access through [role-based access control](#role-based-access-control) or [resource-based access control](#resource-based-access-control).

## Role-based access control

Role-based access control (RBAC) restricts [actions](#action) over [resources](#resource) to authorized users. Access is specified by [roles](#role) assigned to users, using a [manifest file](../HowTo/Use-Manifest-File/Role.md) or an [identity provider](https://auth0.com/docs/authorization/rbac/roles/create-roles).

See the [full list of RBAC permissions](../Reference/RBAC-Permissions.md).

## Resource-based access control

Resource-based access control restricts access to [resources](#resource) to authorized users. Access is specified by allowed [tenants](#tenant) for each resource, using a [manifest file](../HowTo/Use-Manifest-File/Overview.md).

## Terminology

### Action

An action is a functionality of your application to be restricted to authorized users. For example, read, create, sign, encrypt, delete, and destroy.

### Resource

A resource represents a business entity to be managed by your application. Authorization restricts access over resources. QKM currently has the following resources:

| Name | Description |
| :-: | :-: |
| Secret | A key-value element stored in a secure vault system. |
| Key | A cryptographic key. |
| Ethereum account | A cryptographic key allowing interaction with the Ethereum network. |
| [Vault](Stores.md#Vault) | Vault client connector used to persist resources remotely. |
| [Store](Stores.md) | A storage space for a set of secrets, keys, or Ethereum accounts. |
| [Node](Nodes.md) | A representation of an underlying blockchain node. |
| Alias | A representation of an external public key. For example, a [Tessera](https://docs.tessera.consensys.net/en/stable/) address. |
| Registry | A storage space for clarifying a set of aliases |

### Tenant

A tenant is a set of users with the highest access level to [resources](#resource). In [resource-based access control](#resource-based-access-control), you must pass a list of allowed tenants when defining a resource [manifest file](../HowTo/Use-Manifest-File/Overview.md).

### Permission

A permission is an authorization of an [action](#action) over a [resource](#resource), used in [role-based access control (RBAC)](#role-based-access-control). [Permissions](../Reference/RBAC-Permissions.md) take the form `action:resource` and are not mutually exclusive.

### Role

A role is a named set of [permissions](#permission) defined in a [manifest file](../HowTo/Use-Manifest-File/Role.md). Alternatively, you can [use Auth0 to specify roles](https://auth0.com/docs/authorization/rbac/roles/create-roles) and attach permissions to your token.
