package unifi

import (
	"encoding/json"
	"fmt"
)

// ListFirewallPolicies returns all firewall policies.
func (c *Client) ListFirewallPolicies(params ListParams) ([]FirewallPolicy, error) {
	data, err := c.Get(c.siteURL("/firewall/policies") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[FirewallPolicy]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateFirewallPolicy creates a new firewall policy.
func (c *Client) CreateFirewallPolicy(req FirewallPolicyCreateRequest) (*FirewallPolicy, error) {
	data, err := c.Post(c.siteURL("/firewall/policies"), req)
	if err != nil {
		return nil, err
	}
	var policy FirewallPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

// UpdateFirewallPolicy updates a firewall policy (PATCH).
func (c *Client) UpdateFirewallPolicy(id string, req FirewallPolicyCreateRequest) (*FirewallPolicy, error) {
	data, err := c.Patch(c.siteURL(fmt.Sprintf("/firewall/policies/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var policy FirewallPolicy
	if err := json.Unmarshal(data, &policy); err != nil {
		return nil, err
	}
	return &policy, nil
}

// DeleteFirewallPolicy deletes a firewall policy.
func (c *Client) DeleteFirewallPolicy(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/firewall/policies/%s", id)))
	return err
}

// ReorderFirewallPolicies reorders firewall policies.
func (c *Client) ReorderFirewallPolicies(ids []string) error {
	_, err := c.Put(c.siteURL("/firewall/policies/ordering"), OrderingRequest{IDs: ids})
	return err
}

// ListFirewallZones returns all firewall zones.
func (c *Client) ListFirewallZones(params ListParams) ([]FirewallZone, error) {
	data, err := c.Get(c.siteURL("/firewall/zones") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[FirewallZone]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateFirewallZone creates a new firewall zone.
func (c *Client) CreateFirewallZone(req FirewallZoneCreateRequest) (*FirewallZone, error) {
	data, err := c.Post(c.siteURL("/firewall/zones"), req)
	if err != nil {
		return nil, err
	}
	var zone FirewallZone
	if err := json.Unmarshal(data, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// UpdateFirewallZone updates a firewall zone (PUT).
func (c *Client) UpdateFirewallZone(id string, req FirewallZoneCreateRequest) (*FirewallZone, error) {
	data, err := c.Put(c.siteURL(fmt.Sprintf("/firewall/zones/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var zone FirewallZone
	if err := json.Unmarshal(data, &zone); err != nil {
		return nil, err
	}
	return &zone, nil
}

// DeleteFirewallZone deletes a firewall zone.
func (c *Client) DeleteFirewallZone(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/firewall/zones/%s", id)))
	return err
}
