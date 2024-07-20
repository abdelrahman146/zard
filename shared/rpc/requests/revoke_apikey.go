package requests

type RevokeApiKeyRequest struct {
	ApiKey string `json:"api_key" validate:"required"`
}

func (r *RevokeApiKeyRequest) Subject() string {
	return "identity.apikey.revoke"
}

func (r *RevokeApiKeyRequest) Consumer(group string) string {
	return "identity_apikey_revoke_" + group
}

type RevokeApiKeyResponse struct {
	Ok bool `json:"ok"`
}
