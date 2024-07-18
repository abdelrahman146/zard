package auth

//type Claims struct {
//	Iss string    `json:"iss"` // Issuer (Example: "zard")
//	Sub string    `json:"sub"` // Subject (Example: "userId" | "appId")
//	Aud string    `json:"aud"` // Audience (Example: "zard" | "external")
//	Exp time.Time `json:"exp"` // Expiration Time (Example: 1596230000)
//	Iat time.Time `json:"iat"` // Issued At (Example: 1596226400)
//	Jti string    `json:"jti"` // JWT ID (Example: "a-123")
//}

type Claims interface {
	AuthType() string
	User() string
	Permissions() []string
	Organization() string
	Workspace() string
	Subscriptions() []string
}
