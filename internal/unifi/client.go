package unifi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const basePath = "/proxy/network/integration/v1"

// Client is the UniFi Integration API client.
type Client struct {
	BaseURL    string
	APIKey     string
	SiteID     string
	HTTPClient *http.Client
}

// NewClient creates a new UniFi API client.
func NewClient(host, apiKey, siteID string, insecure bool) *Client {
	host = strings.TrimRight(host, "/")
	transport := &http.Transport{}
	if insecure {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	return &Client{
		BaseURL: host,
		APIKey:  apiKey,
		SiteID:  siteID,
		HTTPClient: &http.Client{
			Transport: transport,
		},
	}
}

// ListParams holds common list/pagination parameters.
type ListParams struct {
	Limit  int
	Offset int
	Filter string
}

func (p ListParams) query() string {
	vals := url.Values{}
	if p.Limit > 0 {
		vals.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.Offset > 0 {
		vals.Set("offset", fmt.Sprintf("%d", p.Offset))
	}
	if p.Filter != "" {
		vals.Set("filter", p.Filter)
	}
	if len(vals) == 0 {
		return ""
	}
	return "?" + vals.Encode()
}

func (c *Client) siteURL(path string) string {
	return fmt.Sprintf("%s%s/sites/%s%s", c.BaseURL, basePath, c.SiteID, path)
}

func (c *Client) globalURL(path string) string {
	return fmt.Sprintf("%s%s%s", c.BaseURL, basePath, path)
}

func (c *Client) newRequest(method, url string, body any) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.APIKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func (c *Client) do(req *http.Request) ([]byte, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

// Get performs a GET request and returns raw bytes.
func (c *Client) Get(url string) ([]byte, error) {
	req, err := c.newRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Post performs a POST request.
func (c *Client) Post(url string, body any) ([]byte, error) {
	req, err := c.newRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Put performs a PUT request.
func (c *Client) Put(url string, body any) ([]byte, error) {
	req, err := c.newRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Patch performs a PATCH request.
func (c *Client) Patch(url string, body any) ([]byte, error) {
	req, err := c.newRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}

// Delete performs a DELETE request.
func (c *Client) Delete(url string) ([]byte, error) {
	req, err := c.newRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.do(req)
}
