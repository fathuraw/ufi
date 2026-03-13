package unifi

import (
	"encoding/json"
	"fmt"
)

// ListClients returns all clients for the site.
func (c *Client) ListClients(params ListParams) ([]NetworkClient, error) {
	data, err := c.Get(c.siteURL("/clients") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[NetworkClient]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetClient returns a single client by ID.
func (c *Client) GetClient(id string) (*NetworkClient, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/clients/%s", id)))
	if err != nil {
		return nil, err
	}
	var client NetworkClient
	if err := json.Unmarshal(data, &client); err != nil {
		return nil, err
	}
	return &client, nil
}

// BlockClient blocks a client.
func (c *Client) BlockClient(id string) error {
	_, err := c.Post(c.siteURL(fmt.Sprintf("/clients/%s/actions", id)), ClientAction{Action: "block"})
	return err
}

// UnblockClient unblocks a client.
func (c *Client) UnblockClient(id string) error {
	_, err := c.Post(c.siteURL(fmt.Sprintf("/clients/%s/actions", id)), ClientAction{Action: "unblock"})
	return err
}
