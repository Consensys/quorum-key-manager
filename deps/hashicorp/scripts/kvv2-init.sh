VAULT_SSL_PARAMS=""
if [ -n "$VAULT_CACERT" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --cacert $VAULT_CACERT"
fi  

if [ -n "$VAULT_CLIENT_CERT" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --cert $VAULT_CLIENT_CERT"
fi     

if [ -n "$VAULT_CLIENT_KEY" ]; then
 VAULT_SSL_PARAMS="$VAULT_SSL_PARAMS --key $VAULT_CLIENT_KEY"
fi  

# Enable kv-v2 secret engine
curl --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST ${VAULT_SSL_PARAMS}\
     --data '{"type": "kv-v2", "config": {"force_no_cache": true} }' \
     ${VAULT_ADDR}/v1/sys/mounts/secret

