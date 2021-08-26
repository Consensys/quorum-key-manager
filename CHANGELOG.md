# Quorum Key Manager Release Notes

## v21.7.0-alpha.2 (2021-08-26)
### 🆕 Features
* Support for authorization using OIDC, TLS and API-KEY
* Support for authentication based on roles and permissions
* Usage Postgres DB to resources public information

### 🛠 Bug fixes
* Behaviour alignment over every support key vault
* Keys and secrets were not available after restoring

## v21.7.0-alpha.1 (2021-07-06)
### 🆕 Features
Initial release of the Quorum Key Manager

* Support for [Hashicorp KV Secrets Engine](https://www.vaultproject.io/docs/secrets/kv/kv-v2)
* Support for [Hashicorp keys plugin](https://github.com/ConsenSys/orchestrate-hashicorp-vault-plugin) (custom plugin)
* Support for [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/) (secrets and keys)
* Support for [AWS Secrets Manager](https://aws.amazon.com/secrets-manager/)
* Support for [AWS Key Management Service](https://aws.amazon.com/kms/) (KMS)
* Support for Ethereum account management using an underlying key store 
* Node proxy connected to an underlying Blockchain Node (tested with [GoQuorum](https://docs.goquorum.consensys.net/en/stable/) and [Hyperledger Besu](https://www.hyperledger.org/use/besu)) intercepting JSON-RPC calls
