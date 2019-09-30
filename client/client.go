package client

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	// TokenURL is the URL used to get a token so we can access our application
	// setup in the Spotify dashboard.
	TokenURL = "https://accounts.spotify.com/api/token"
	BaseURL  = "https://api.spotify.com/v1"
)

// Client holds the http client used to make requests, the base url and the user agent.
type Client struct {
	ClientSecret string
	ClientID     string
	Token        string
	BearerToken  string

	HttpClient *http.Client
}

func NewClient() Client {
	return Client{
		ClientSecret: os.Getenv("SECRET"),
		ClientID:     os.Getenv("CLIENT_ID"),
		HttpClient:   &http.Client{},
	}
}

// Auth authorises us with the Spotify API.
func (c *Client) Auth() error {
	var authResponse AuthResponse

	form := url.Values{}
	form.Add("grant_type", "client_credentials")

	req, err := c.NewAuthRequest("POST", strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}

	authKey := fmt.Sprintf("%s:%s", c.ClientID, c.ClientSecret)
	encoded := base64.StdEncoding.EncodeToString([]byte(authKey))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", encoded))

	_, err = c.Do(req, &authResponse)
	if err != nil {
		return err
	}

	c.BearerToken = authResponse.AccessToken

	return nil
}

// Do performs the API request.
func (c *Client) Do(req *http.Request, value interface{}) (*http.Response, error) {
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Printf("resp %+v", resp)
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(value)

	return resp, err
}

// NewAuthRequest creates a new auth request to the Spotify API so we can get a bearer token.
func (c *Client) NewAuthRequest(method string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, TokenURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req, nil
}

// NewRequest creates a new request based on a method, path and then a response body to marshall the data into.
func (c *Client) NewRequest(method string, path string, body interface{}) (*http.Request, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, errors.Wrap(
			err, "Problem marshalling payload for spotify request",
		)
	}

	req, err := http.NewRequest(method, c.generateEndpoint(path), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c Client) generateEndpoint(path string) string {
	return fmt.Sprintf(
		"%s/%s",
		BaseURL,
		path,
	)
}
