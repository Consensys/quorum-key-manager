---
title: Connect to an Infura endpoint
description: Connect to an Infura endpoint
sidebar_position: 3
---

# Connect to an Infura endpoint

This tutorial walks you through connecting Quorum Key Manager (QKM) to an Infura endpoint.

This tutorial demonstrates connecting to Infura to allow a QKM-managed account to interact with a smart contract on the Rinkeby network. It uses a QKM instance on top of a Kubernetes cluster with an AWS KMS as key storage.

## Prerequisites

- An [Infura](https://infura.io/) account
- Access an [Amazon EKS cluster](https://docs.aws.amazon.com/eks/latest/userguide/clusters.html) or [minikube](https://minikube.sigs.k8s.io/docs/start/)
- The following tools ready and aligned with your cluster Kubernetes version:
  - [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
  - [Helm](https://helm.sh/)
  - [Helmfile](https://github.com/roboll/helmfile)
- An [AWS key manager](https://aws.amazon.com/kms/) ready with the following associated credentials:
  - `ACCESS_ID`
  - `SECRET_KEY`
  - `REGION`

## Steps

1. In Infura, create a new Ethereum project.

1. Go to the project settings, and get the endpoint associated with your target network (in this example, Rinkeby). This endpoint is your node URL.

   ```text
   https://rinkeby.infura.io/v3/<YOUR_PROJECT_ID>
   ```

1. Start QKM with the [Helmfile chart](https://github.com/ConsenSys/quorum-key-manager-kubernetes) specifying a manifest file with the following content:

   ```yml
   # Infura node manifest
   - kind: Node
     name: infura-node
     specs:
       rpc:
         addr: https://rinkeby.infura.io/v3/<YOUR_PROJECT_ID>

   # Ethereum store manifest backed by an AWS keystore
   - kind: Vault
     type: aws
     name: aws-europe
     specs:
       access_id: <YOUR_KMS_ACCOUNT_ACCESS_ID>
       secret_key: <YOUR_KMS_ACCOUNT_SECRET>
       region: <YOUR_KMS_ACCOUNT_REGION>
   - kind: Store
     type: key
     name: aws-keys
     specs:
       vault: aws-europe
   - kind: Store
     type: ethereum
     name: eth-accounts
     specs:
       key_store: aws-keys
   ```

1. You can connect QKM to Infura using one of the following methods:

   - Port forwarding. Run the following commands:

     ```bash
     export POD_NAME=$(kubectl get pods --namespace $QKM-NAMESPACE -l "app.kubernetes.io/name=quorumkeymanager,app.kubernetes.io/instance=quorum-key-manager" -o jsonpath="{.items[0].metadata.name}")

     kubectl --namespace $QKM-NAMESPACE port-forward $POD_NAME 8080:8080
     ```

   - [Ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/). Set up the appropriate Ingress values in the helm chart before deployment in order to have Ingress activated and configured according to your solution.

1. Create a new Ethereum account using the `createAccount` method, filling in your own values. Change the `localhost` target according to your Ingress configuration and required URL:

   ```bash
   curl --location --request POST 'https://localhost:8080/stores/eth-accounts/ethereum' \
       --header 'Authorization: Basic YWRtaW4tdXNlcg==' \
       --header 'Content-Type: application/json' \
       --data-raw '{
           "keyId": "my-infura-key",
           "tags": {
               "owner": "mySelf"
           }
       }'
   ```

   The response yields an Ethereum account address.

1. Send JSON-RPC-based transactions to the Infura node created in step 1 using the Ethereum account created in step 5. Change `<MY_ETH_ACCOUNT_ADDRESS>` to the response address from the previous step, and the `localhost` target according to your Ingress configuration and required URL:

   ```bash
   curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_sendTransaction","params":[{"from": <MY_ETH_ACCOUNT_ADDRESS>,"to": "0x015C7C7A7D65bbdb117C573007219107BD7486f9","value": "0x1000000"}], "id":1}' http://localhost:8080/nodes/rinkeby-infura
   ```

1. Test that everything worked:

   - Check that Infura has recorded your action with `<YOUR_PROJECT_ID>`.
   - Check that the transaction has been mined. You can view it on Etherscan.

## Next steps

Once connected to an Infura endpoint, you can [send Ethereum meta-transactions](SendMetaTxn.md).
