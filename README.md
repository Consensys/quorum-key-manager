# Quorum Key Manager

## License

Orchestrate is licensed under the BSL 1.1.

Please refer to the [LICENSE file](LICENSE) for a detailed description of the license.

Please contact [orchestrate@consensys.net](mailto:orchestrate@consensys.net) if you need to purchase a license for a production use-case.  



## Aws specific annotations

When getting one Key, KM may return the following annotations :

- "aws-KeyID" : the Key inner ID
- "aws-KeyStoreID" : the Key's Keystore ID
- "aws-ClusterHSMID" : the Key's Cluster HSM ID, when backed by an actual HSM
- "aws-AccountID" : the current user (requester) account ID
- "awsARN" : the key Amazon Ressource Name