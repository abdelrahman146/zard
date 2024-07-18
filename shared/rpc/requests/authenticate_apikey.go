package requests

type AuthenticateApiKeyRequest struct {
	ApiKey string `json:"api_key" validate:"required"`
}

func (r *AuthenticateApiKeyRequest) Subject() string {
	return "identity.apikey.authenticate"
}

func (r *AuthenticateApiKeyRequest) Consumer(group string) string {
	return "identity_apikey_authenticate_" + group
}

type AuthenticateApiKeyResponse struct {
	Org  string   `json:"org"`
	Ws   string   `json:"workspace"`
	Subs []string `json:"sub"`
}

func (r *AuthenticateApiKeyResponse) AuthType() string {
	return "apikey"
}

func (r *AuthenticateApiKeyResponse) User() string {
	return ""
}

func (r *AuthenticateApiKeyResponse) Permissions() []string {
	return []string{}
}

func (r *AuthenticateApiKeyResponse) Organization() string {
	return r.Org
}

func (r *AuthenticateApiKeyResponse) Workspace() string {
	return r.Ws
}

func (r *AuthenticateApiKeyResponse) Subscriptions() []string {
	return r.Subs
}
