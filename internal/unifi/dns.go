package unifi

import (
	"encoding/json"
	"fmt"
)

// ListDNSRecords returns all DNS records.
func (c *Client) ListDNSRecords(params ListParams) ([]DNSRecord, error) {
	data, err := c.Get(c.siteURL("/dns/policies") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[DNSRecord]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// CreateDNSRecord creates a new DNS record.
func (c *Client) CreateDNSRecord(req DNSRecordCreateRequest) (*DNSRecord, error) {
	data, err := c.Post(c.siteURL("/dns/policies"), req)
	if err != nil {
		return nil, err
	}
	var record DNSRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// UpdateDNSRecord updates a DNS record.
func (c *Client) UpdateDNSRecord(id string, req DNSRecordCreateRequest) (*DNSRecord, error) {
	data, err := c.Put(c.siteURL(fmt.Sprintf("/dns/records/%s", id)), req)
	if err != nil {
		return nil, err
	}
	var record DNSRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, err
	}
	return &record, nil
}

// DeleteDNSRecord deletes a DNS record.
func (c *Client) DeleteDNSRecord(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/dns/records/%s", id)))
	return err
}
