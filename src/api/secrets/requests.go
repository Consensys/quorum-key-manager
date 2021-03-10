package secrets

type CreateSecretRequest struct {
	ID    string            `json:"id" example:"my-privateKey" validate:"required"`
	Value string            `json:"value" example:"fa88c4a5912f80503d6b5503880d0745f4b88a1ff90ce8f64cdd8f32cc3bc249" validate:"required"`
	Tags  map[string]string `json:"tags"`
}
