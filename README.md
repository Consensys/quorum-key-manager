# Quorum Key Manager


## Aws specific annotations

When getting one Key, KM may return the following annotations :

- "aws-KeyID" : the Key inner ID
- "aws-KeyStoreID" : the Key's Keystore ID
- "aws-ClusterHSMID" : the Key's Cluster HSM ID, when backed by an actual HSM
- "aws-AccountID" : the current user (requester) account ID
- "awsARN" : the key Amazon Ressource Name