---
title: Use TLS
description: How to authenticate QKM using TLS.
sidebar_position: 2
---

# Authenticate using TLS

You can [authenticate](../../Concepts/Authentication.md#authentication) incoming Quorum Key Manager (QKM) requests using mutual TLS authentication.

To use TLS mutual authentication, start QKM in SSL mode with the [`--https-enable`](../../Reference/CLI/CLI-Syntax.md#https-enable), [`--https-server-cert`](../../Reference/CLI/CLI-Syntax.md#https-server-cert), and [`--https-server-key`](../../Reference/CLI/CLI-Syntax.md#https-server-key) command line options, and specify a TLS certificate authority (CA) certificate with the [`--auth-tls-ca`](../../Reference/CLI/CLI-Syntax.md#auth-tls-ca) option.

:::info

Starting Quorum Key Manager with TLS authentication

```bash
key-manager run --https-enable --https-server-cert=tls.crt --https-server-key=tls.key --auth-tls-ca=ca.crt --manifest-path=/config/default.yml
```

:::

## TLS certificate

The CA certificate must contain one or more CAs to validate client certificates presented to QKM.

If a client presents a valid certificate signed by one of the CAs, then the client is authenticated.

QKM extracts the following user information from the subject field of the client certificate:

- Username and optional [tenant](../../Concepts/Authorization.md#tenant) from the common name (CN) (for example, `/CN=tenant|user` or `/CN=user`)
- [Roles](../../Concepts/Authorization.md#role) from the certificate's organization (O) (for example, `/O=role1/O=role2`)
- [Permissions](../../Concepts/Authorization.md#permission) from the certificate's organization unit (OU) (for example, `/OU=*:read/OU=secret:write`)

You can use the `openssl` command line tool to generate a certificate signing request:

:::info

Example certificate signing request

```bash
openssl req -new -key jbeda.pem -out jbeda-csr.pem -subj "/CN=auth0|alice/O=admin/OU=sign:eth1Account"
```

:::
