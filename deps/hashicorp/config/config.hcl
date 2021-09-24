backend "file" {
  path = "/vault/file"
}

listener "tcp" {
  address = "hashicorp:8200"
  tls_cert_file = "/vault/tls/vault.crt"
  tls_client_ca_file = "/vault/tls/root.crt"
  tls_key_file = "/vault/tls/vault.key"
}

default_lease_ttl = "15m"
max_lease_ttl = 99999999
api_addr = "http://hashicorp:8200"
plugin_directory = "/vault/plugins"
log_level = "Debug"

ui = false
