package unifi

import (
	"encoding/json"
	"fmt"
)

// ListDevices returns all devices for the site.
func (c *Client) ListDevices(params ListParams) ([]Device, error) {
	data, err := c.Get(c.siteURL("/devices") + params.query())
	if err != nil {
		return nil, err
	}
	var resp ListResponse[Device]
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// GetDevice returns a single device by ID (direct object, not wrapped).
func (c *Client) GetDevice(id string) (*Device, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/devices/%s", id)))
	if err != nil {
		return nil, err
	}
	var device Device
	if err := json.Unmarshal(data, &device); err != nil {
		return nil, err
	}
	return &device, nil
}

// AdoptDevice adopts a device by MAC address.
func (c *Client) AdoptDevice(mac string) error {
	_, err := c.Post(c.siteURL("/devices"), AdoptRequest{MAC: mac})
	return err
}

// RemoveDevice removes a device by ID.
func (c *Client) RemoveDevice(id string) error {
	_, err := c.Delete(c.siteURL(fmt.Sprintf("/devices/%s", id)))
	return err
}

// RestartDevice restarts a device.
func (c *Client) RestartDevice(id string) error {
	_, err := c.Post(c.siteURL(fmt.Sprintf("/devices/%s/actions", id)), DeviceAction{Action: "restart"})
	return err
}

// GetDeviceStatistics returns the latest statistics for a device (direct object, not wrapped).
func (c *Client) GetDeviceStatistics(id string) (*DeviceStatistics, error) {
	data, err := c.Get(c.siteURL(fmt.Sprintf("/devices/%s/statistics/latest", id)))
	if err != nil {
		return nil, err
	}
	var stats DeviceStatistics
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, err
	}
	return &stats, nil
}

// PowerCyclePort power cycles a port on a device.
func (c *Client) PowerCyclePort(deviceID string, portIdx int) error {
	_, err := c.Post(
		c.siteURL(fmt.Sprintf("/devices/%s/ports/%d/actions", deviceID, portIdx)),
		PortAction{Action: "power_cycle"},
	)
	return err
}
