package unifi

import (
	"encoding/json"
	"fmt"
)

// ListNetworks returns all networks for the site.
func (c *Client) ListNetworks(params ListParams) ([]Network, error) {
	data, err := c.Get(c.siteURL("/networks") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Network]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetNetwork returns a single network by ID.
func (c *Client) GetNetwork(id string) (*Network, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/networks/%s", id)))
	if err != nil {
		return nil, err
	}
	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

// CreateNetwork creates a new network.
func (c *Client) CreateNetwork(req NetworkCreateRequest) (*Network, error) {
	data, err := c.Post(c.siteURL("/networks"), req)
	if err != nil {
		return nil, err
	}
	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

// UpdateNetwork updates an existing network.
func (c *Client) UpdateNetwork(id string, req NetworkCreateRequest) (*Network, error) {
	data, err := c.Put(c.siteURL(fmt.Sprintf("/networks/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var network Network
	if err := json.Unmarshal(data, &network); err != nil {
		return nil, err
	}
	return &network, nil
}

// DeleteNetwork deletes a network by ID.
func (c *Client) DeleteNetwork(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/networks/%s", id)))
	return err
}

// GetNetworkReferences returns dependency references for a network.
func (c *Client) GetNetworkReferences(id string) ([]NetworkReference, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/networks/%s/references", id)))
	if err != nil {
		return nil, err
	}
	var resp ListResponse[NetworkReference]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}
