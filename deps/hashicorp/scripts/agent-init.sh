VAULT_TOKEN=$(cat "${ROOT_TOKEN_PATH}")

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
 
curl -s --header "X-Vault-Token: ${VAULT_TOKEN}" --request POST ${VAULT_SSL_PARAMS} \
  --data '{"type": "approle"}' \
  ${VAULT_ADDR}/v1/sys/auth/approle

curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request PUT ${VAULT_SSL_PARAMS} \
  --data '{ "policy":"path \"'"${PLUGIN_MOUNT_PATH}/*"'\" { capabilities = [\"create\", \"read\", \"update\", \"delete\", \"list\"] }" }' \
  ${VAULT_ADDR}/v1/sys/policies/acl/allow_secrets

curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request POST ${VAULT_SSL_PARAMS} \
  --data '{"policies": ["allow_secrets"]}' \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager

curl -s --header "X-Vault-Token: $VAULT_TOKEN" ${VAULT_SSL_PARAMS} \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager/role-id > role.json
ROLE_ID=$(cat role.json | jq .data.role_id | tr -d '"')
echo $ROLE_ID > ${ROLE_FILE_PATH}
rm role.json
  
curl -s --header "X-Vault-Token: $VAULT_TOKEN" --request POST ${VAULT_SSL_PARAMS} \
  ${VAULT_ADDR}/v1/auth/approle/role/key-manager/secret-id > secret.json
SECRET_ID=$(cat secret.json | jq .data.secret_id | tr -d '"')
echo $SECRET_ID > ${SECRET_FILE_PATH}
rm secret.json
