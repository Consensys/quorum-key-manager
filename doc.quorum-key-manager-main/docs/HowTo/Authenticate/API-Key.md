---
title: Use API keys
description: How to authenticate QKM using an API key.
sidebar_position: 3
---

# Authenticate using API keys

You can [authenticate](../../Concepts/Authentication.md#authentication) incoming Quorum Key Manager (QKM) requests using API keys.

Specify an API key file with the [`--auth-api-key-file`](../../Reference/CLI/CLI-Syntax.md#auth-api-key-file) command line option when starting QKM.

:::info

Starting Quorum Key Manager with API key authentication

```bash
key-manager run --auth-api-key-file=api_key_file.csv --manifest-path=/config/default.yml
```

:::

## API key file

The API key file is a CSV file with four columns:

- sha256({apiKey})
- username and optional [tenant](../../Concepts/Authorization.md#tenant)
- [permissions](../../Reference/RBAC-Permissions.md)
- [roles](../../Concepts/Authorization.md#role)

Each CSV line must be a unique API key and all API keys must be in UUID V4 format.

:::info Example API key file

```
sha256({apiKey1}),tenant1|username1,"*:secret,*:keys","role-admin"
sha256({apiKey2}),username2,"read:*","role-guest"
```

:::

To extract an API key, QKM uses the standard [HTTP basic authentication](https://swagger.io/docs/specification/authentication/basic-authentication/) scheme with a blank username and the API key as the password:

<!--tabs-->

# Syntax

```
Authorization: Basic <base64({apiKey})>
```

# Example

```
Authorization: Basic OjA2ZGExYWZlLTE2ZDMtNDhmZS04ZWMyLWZlYTg2NDhkNzM3YQ==
```

<!--/tabs-->

If a user passes an API key that's in the CSV file, user information from the corresponding line in the CSV file is attached to the request.
