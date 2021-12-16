# Quorum Key Manager Release Notes

## v21.12.0 LTS (2021-12-14)
### 🆕 Features
* Support for import of secrets, keys and ethereum accounts through command line (`sync` command)
* Support for alias management on `/registries/{registryName}/aliases`
* Support Token Issuer Servers to validate JWTs. Environment variable `AUTH_JWT_ISSUER_URL` and `AUTH_JWT_AUDIENCE`

### 🛠 Changes
* Env var `AUTH_OIDC_CA_CERT` and flag `--auth-oidc-ca-cert` renamed to `AUTH_OIDC_PUB_KEY` and `--auth-oidc-pub-key`
* Manifest definition changes introducing the new kind `Vault`. See the documentation for more information.
* Removed usage of `AUTH_JWT_CERTIFICATE` in favor of `AUTH_JWT_ISSUER_URL` and `AUTH_JWT_AUDIENCE`

## v21.9.3 (2021-11-10)
### 🛠 Bug fixes
* Fixes bug in Hashicorp client that prevents the process from exiting when a new token is written or updated from filesystem.

## v21.9.2 (2021-10-18)
### 🛠 Bug fixes
* Use comma as column separator in CSV file for API key definition

## v21.9.1 (2021-10-05)
### 🛠 Bug fixes
* Enabled support for TLS communication with Hashicorp Vault

## v21.9.0 (2021-09-22)
### 🆕 Features
Initial release of the Quorum Key Manager

* Support for [Hashicorp KV Secrets Engine](https://www.vaultproject.io/docs/secrets/kv/kv-v2)
* Support for [Quorum Hashicorp Vault Plugin](https://github.com/ConsenSys/quorum-hashicorp-vault-plugin) (custom plugin)
* Support for [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/) (secrets and keys)
* Support for [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
* Support for [AWS Key Management Service](https://aws.amazon.com/kms/) (KMS)
* Support for Ethereum account management using an underlying key store 
* Node proxy connected to an underlying Blockchain Node (tested with [GoQuorum](https://docs.goquorum.consensys.net/en/stable/) and [Hyperledger Besu](https://www.hyperledger.org/use/besu)) intercepting JSON-RPC calls
* Support for authorization using OIDC, TLS and API-KEY
* Support for authentication based on roles and permissions
* Usage Postgres DB to resources public information
* Support for PostgreSQL migrations through command line

### 🐛 Know issues
* Communication between HashiCorp Vault and Quorum Key Manager cannot use TLS
