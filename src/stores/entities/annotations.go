package entities

type Annotation struct {
	AWSKeyID             string `json:"AWSKeyID,omitempty"`
	AWSCustomKeyStoreID  string `json:"AWSCustomKeyStoreID,omitempty"`
	AWSCloudHsmClusterID string `json:"AWSCloudHsmClusterID,omitempty"`
	AWSAccountID         string `json:"AWSAccountID,omitempty"`
	AWSArn               string `json:"AWSArn,omitempty"`
}
