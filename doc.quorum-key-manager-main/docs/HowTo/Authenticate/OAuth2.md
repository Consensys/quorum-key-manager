---
title: Use OAuth 2.0
description: How to authenticate QKM with OAuth 2.0 and JWTs.
sidebar_position: 1
---

# Authenticate using OpenID Connect

You can [authenticate](../../Concepts/Authentication.md#authentication) incoming Quorum Key Manager (QKM) requests with the [OAuth 2.0](https://oauth.net/2/) standard using [JSON Web Tokens (JWTs)](https://jwt.io/).

To use OAuth 2.0 authentication, user requests must present a JWT through the HTTP `Authorization` header with value `Bearer <token>`.

Refer to the [OAuth 2.0](https://oauth.net/2/) and [OpenID Connect (OIDC)](https://openid.net/specs/openid-connect-core-1_0.html) documentation for detailed information.

## Command line options

You can set the following options at QKM runtime to configure OAuth 2.0 authentication.

- [`--auth-oidc-issuer-url`](../../Reference/CLI/CLI-Syntax.md#auth-oidc-issuer-url) - URL of the OpenID Connect server.
- [`--auth-oidc-audience`](../../Reference/CLI/CLI-Syntax.md#auth-oidc-audience) - Expected audience in access tokens.

:::info

Starting Quorum Key Manager with OAuth 2.0 authentication

```bash
key-manager run --auth-oidc-issuer-url="https://quorum-key-manager.eu.auth0.com" --auth-oidc-audience=https://quorum-key-manager.consensys.net --manifest-path=/config/default.yml
```

:::
