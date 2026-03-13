package unifi

import (
	"encoding/json"
	"fmt"
)

// ListWiFiBroadcasts returns all WiFi broadcasts.
func (c *Client) ListWiFiBroadcasts(params ListParams) ([]WiFiBroadcast, error) {
	data, err := c.Get(c.siteURL("/wifi/broadcasts") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[WiFiBroadcast]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetWiFiBroadcast returns a single WiFi broadcast.
func (c *Client) GetWiFiBroadcast(id string) (*WiFiBroadcast, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/wifi/broadcasts/%s", id)))
	if err != nil {
		return nil, err
	}
	var broadcast WiFiBroadcast
	if err := json.Unmarshal(data, &broadcast); err != nil {
		return nil, err
	}
	return &broadcast, nil
}

// CreateWiFiBroadcast creates a new WiFi broadcast.
func (c *Client) CreateWiFiBroadcast(req WiFiCreateRequest) (*WiFiBroadcast, error) {
	data, err := c.Post(c.siteURL("/wifi/broadcasts"), req)
	if err != nil {
		return nil, err
	}
	var broadcast WiFiBroadcast
	if err := json.Unmarshal(data, &broadcast); err != nil {
		return nil, err
	}
	return &broadcast, nil
}

// UpdateWiFiBroadcast updates a WiFi broadcast.
func (c *Client) UpdateWiFiBroadcast(id string, req WiFiCreateRequest) (*WiFiBroadcast, error) {
	data, err := c.Put(c.siteURL(fmt.Sprintf("/wifi/broadcasts/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var broadcast WiFiBroadcast
	if err := json.Unmarshal(data, &broadcast); err != nil {
		return nil, err
	}
	return &broadcast, nil
}

// DeleteWiFiBroadcast deletes a WiFi broadcast.
func (c *Client) DeleteWiFiBroadcast(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/wifi/broadcasts/%s", id)))
	return err
}
