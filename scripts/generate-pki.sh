#!/bin/bash

# get cfssl + cfssljson
go get github.com/cloudflare/cfssl/cmd/cfssl
go get github.com/cloudflare/cfssl/cmd/cfssljson

# useful dirs
CONF_DIR=./deps/cfssl
GEN_DIR=./deps/cfssl/generated
DEST_DIR_CA=./deps/config/ca
DEST_DIR_CERTS=./deps/config/certificates
DEST_DIR_VAULT_CERTS=./deps/hashicorp/tls
DEST_DIR_PG_CERTS=./deps/postgres/tls

mkdir -p $GEN_DIR $DEST_DIR_CA $DEST_DIR_CERTS $DEST_DIR_VAULT_CERTS $DEST_DIR_PG_CERTS

ROOT_CRT=./deps/cfssl/generated/ca.pem
ROOT_KEY=./deps/cfssl/generated/ca-key.pem

# Gen root + intermediate
cfssl gencert -initca $CONF_DIR/root.json | cfssljson -bare $GEN_DIR/ca
cfssl gencert -initca $CONF_DIR/intermediate.json | cfssljson -bare $GEN_DIR/intermediate_ca
cfssl sign -ca $ROOT_CRT -ca-key $ROOT_KEY -config $CONF_DIR//cfssl.json -profile intermediate_ca $GEN_DIR/intermediate_ca.csr | cfssljson -bare $GEN_DIR/intermediate_ca

INTER_CRT=./deps/cfssl/generated/intermediate_ca.pem
INTER_KEY=./deps/cfssl/generated/intermediate_ca-key.pem

# Gen leaves from intermediate
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=client $CONF_DIR/qkm-client-no-auth.json | cfssljson -bare $GEN_DIR/qkm-client-no-auth
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=client $CONF_DIR/qkm-client-auth.json | cfssljson -bare $GEN_DIR/qkm-client-auth
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=server $CONF_DIR/qkm-client-auth.json | cfssljson -bare $GEN_DIR/qkm-server
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=server $CONF_DIR/vault.json | cfssljson -bare $GEN_DIR/vault-server
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=client $CONF_DIR/vault.json | cfssljson -bare $GEN_DIR/vault-client
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=server $CONF_DIR/postgres.json | cfssljson -bare $GEN_DIR/postgres-server
cfssl gencert -ca $INTER_CRT -ca-key $INTER_KEY -config $CONF_DIR/cfssl.json -profile=client $CONF_DIR/postgres.json | cfssljson -bare $GEN_DIR/postgres-client

# ca.crt is ca.pem >> intermediate.pem
cat $GEN_DIR/ca.pem > $GEN_DIR/ca.crt
cat $GEN_DIR/intermediate_ca.pem >> $GEN_DIR/ca.crt

# Verify certs
openssl verify -CAfile $GEN_DIR/ca.crt $GEN_DIR/qkm-client-no-auth.pem $GEN_DIR/qkm-client-auth.pem $GEN_DIR/vault-server.pem $GEN_DIR/vault-client.pem $GEN_DIR/qkm-server.pem $GEN_DIR/postgres-server.pem $GEN_DIR/postgres-client.pem

# Relocate
mv $GEN_DIR/ca.crt $DEST_DIR_CA/ca.crt
cp $DEST_DIR_CA/ca.crt $DEST_DIR_VAULT_CERTS/ca.crt
mv $GEN_DIR/qkm-client-auth.pem $DEST_DIR_CERTS/client.crt
mv $GEN_DIR/qkm-client-auth-key.pem $DEST_DIR_CERTS/client.key
mv $GEN_DIR/qkm-client-no-auth.pem $DEST_DIR_CERTS/client_no_auth.crt
mv $GEN_DIR/qkm-client-no-auth-key.pem $DEST_DIR_CERTS/client_no_auth.key
mv $GEN_DIR/qkm-server.pem $DEST_DIR_CERTS/https.crt
mv $GEN_DIR/qkm-server-key.pem $DEST_DIR_CERTS/https.key
mv $GEN_DIR/postgres-server.pem $DEST_DIR_PG_CERTS/tls.crt
mv $GEN_DIR/postgres-server-key.pem $DEST_DIR_PG_CERTS/tls.key
cp $GEN_DIR/intermediate_ca.pem $DEST_DIR_PG_CERTS/ca.crt
mv $GEN_DIR/vault-server.pem $DEST_DIR_VAULT_CERTS/tls.crt
mv $GEN_DIR/vault-server-key.pem $DEST_DIR_VAULT_CERTS/tls.key
mv $GEN_DIR/vault-client.pem $DEST_DIR_VAULT_CERTS/client.crt
mv $GEN_DIR/vault-client-key.pem $DEST_DIR_VAULT_CERTS/client.key

# cleanup
rm -rf $GEN_DIR/*


