package client

import (
	"net/http"
)

// APIClient represents a client for interacting with a hypothetical API.
type APIClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewAPIClient creates a new instance of the API client.
func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}


