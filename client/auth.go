package client

type (
	// AuthResponse is the response from the Spotify API which we use to authenticate API reqursts.
	AuthResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Expiry      int    `json:"expires_in"`
	}
)
