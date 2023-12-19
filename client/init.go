package client

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strings"
)

// APIClient represents an API client
type APIClient struct {
    BaseURL    string
    HTTPClient *http.Client
}

// NewAPIClient creates a new instance of APIClient
func NewAPIClient(baseURL string) *APIClient {
    return &APIClient{
        BaseURL:    baseURL,
        HTTPClient: &http.Client{},
    }
}

// MakeRequest performs an HTTP request to the specified endpoint with the given method, payload, and headers
func (c *APIClient) MakeRequest(endpoint, method string, payload []byte, headers map[string]string) ([]byte, error) {
    url := c.BaseURL + endpoint

    req, err := http.NewRequest(method, url, strings.NewReader(string(payload)))
    if err != nil {
        return nil, fmt.Errorf("error creating request: %v", err)
    }

    for key, value := range headers {
        req.Header.Add(key, value)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("error sending request: %v", err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("error reading response: %v", err)
    }

    return body, nil
}
