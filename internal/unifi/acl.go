package unifi

import (
	"encoding/json"
	"fmt"
)

// ListACLRules returns all ACL rules.
func (c *Client) ListACLRules(params ListParams) ([]ACLRule, error) {
	data, err := c.Get(c.siteURL("/acl-rules") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[ACLRule]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateACLRule creates a new ACL rule.
func (c *Client) CreateACLRule(req ACLRuleCreateRequest) (*ACLRule, error) {
	data, err := c.Post(c.siteURL("/acl-rules"), req)
	if err != nil {
		return nil, err
	}
	var rule ACLRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// UpdateACLRule updates an ACL rule.
func (c *Client) UpdateACLRule(id string, req ACLRuleCreateRequest) (*ACLRule, error) {
	data, err := c.Put(c.siteURL(fmt.Sprintf("/acl-rules/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var rule ACLRule
	if err := json.Unmarshal(data, &rule); err != nil {
		return nil, err
	}
	return &rule, nil
}

// DeleteACLRule deletes an ACL rule.
func (c *Client) DeleteACLRule(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/acl-rules/%s", id)))
	return err
}

// ReorderACLRules reorders ACL rules.
func (c *Client) ReorderACLRules(ids []string) error {
	_, err := c.Put(c.siteURL("/acl-rules/ordering"), OrderingRequest{IDs: ids})
	return err
}
