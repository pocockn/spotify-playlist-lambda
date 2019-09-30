package client

type (
	AuthResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Expiry      int    `json:"expires_in"`
	}
)
