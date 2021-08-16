package entities

type Annotation struct {
	AWSKeyID             string `json:"aws_key_id,omitempty"`
	AWSCustomKeyStoreID  string `json:"aws_custom_key_store_id,omitempty"`
	AWSCloudHsmClusterID string `json:"aws_cloud_hsm_cluster_id,omitempty"`
	AWSAccountID         string `json:"aws_account_id,omitempty"`
	AWSArn               string `json:"aws_arn,omitempty"`
}
