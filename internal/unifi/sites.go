package unifi

import "encoding/json"

// ListSites returns all sites.
func (c *Client) ListSites() ([]Site, error) {
	data, err := c.Get(c.globalURL("/sites"))
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Site]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

