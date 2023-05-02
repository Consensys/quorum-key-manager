---
title: Add a role
description: Add Roles using manifest
sidebar_position: 4
---

# Add a role to Quorum Key Manager

You can define a [role](../../Concepts/Authorization.md#role) in a Quorum Key Manager (QKM) [manifest file](Overview.md).

Use the following fields to configure one or more roles:

- `kind`: _string_ - the string `Role`
- `name`: _string_ - name of the role
- `specs`: _object_ - configuration object containing a list of [permissions](../../Reference/RBAC-Permissions.md) assigned to the role.

```yaml title="Example role manifest file"
# Anonymous role manifest (specifies permissions for anonymous requests)
- kind: Role
  name: anonymous
  specs:
    permissions:
      - "read:nodes"

# Guest role manifest
- kind: Role
  name: guest
  specs:
    permissions:
      - "read:*"

# Signer role manifest
- kind: Role
  name: signer
  specs:
    permissions:
      - "read:*"
      - "sign:keys"
      - "sign:ethereum"

# Admin role manifest
- kind: Role
  name: admin
  specs:
    permissions:
      - "*:*"
```
