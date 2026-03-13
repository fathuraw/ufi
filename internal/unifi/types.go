package unifi

// ListResponse wraps paginated list responses from the UniFi Integration API.
type ListResponse[T any] struct {
	Data       []T  `json:"data"`
	Offset     int  `json:"offset"`
	Limit      int  `json:"limit"`
	TotalCount int  `json:"totalCount"`
	Count      int  `json:"count"`
	HasNext    bool `json:"hasNext"`
}

// Site represents a UniFi site.
type Site struct {
	ID                string `json:"id"`
	InternalReference string `json:"internalReference"`
	Name              string `json:"name"`
}

// Device represents a UniFi device.
type Device struct {
	ID                string   `json:"id"`
	MACAddress        string   `json:"macAddress"`
	IPAddress         string   `json:"ipAddress"`
	Name              string   `json:"name"`
	Model             string   `json:"model"`
	State             string   `json:"state"`
	Supported         bool     `json:"supported"`
	FirmwareVersion   string   `json:"firmwareVersion"`
	FirmwareUpdatable bool     `json:"firmwareUpdatable"`
	AdoptedAt         string   `json:"adoptedAt"`
	ProvisionedAt     string   `json:"provisionedAt"`
	ConfigurationID   string   `json:"configurationId"`
	FeaturesList      []string `json:"features,omitempty"`     // list form (from list endpoint)
	InterfacesList    []string `json:"interfaces,omitempty"`   // list form (from list endpoint)
	Uplink            *struct {
		DeviceID string `json:"deviceId"`
	} `json:"uplink,omitempty"`
	Features *struct {
		Switching   *struct{} `json:"switching,omitempty"`
		AccessPoint *struct{} `json:"accessPoint,omitempty"`
	} `json:"features,omitempty"`
	Interfaces *DeviceInterfaces `json:"interfaces,omitempty"`
}

// DeviceInterfaces holds port and radio info from the detail endpoint.
type DeviceInterfaces struct {
	Ports  []DevicePort  `json:"ports,omitempty"`
	Radios []DeviceRadio `json:"radios,omitempty"`
}

// DevicePort represents a switch port.
type DevicePort struct {
	Idx           int    `json:"idx"`
	State         string `json:"state"`
	Connector     string `json:"connector"`
	MaxSpeedMbps  int    `json:"maxSpeedMbps"`
	SpeedMbps     int    `json:"speedMbps"`
	PoE           *struct {
		Standard string `json:"standard"`
		Enabled  bool   `json:"enabled"`
		State    string `json:"state"`
	} `json:"poe,omitempty"`
}

// DeviceRadio represents a radio interface.
type DeviceRadio struct {
	WLANStandard    string  `json:"wlanStandard"`
	FrequencyGHz    float64 `json:"frequencyGHz"`
	ChannelWidthMHz int     `json:"channelWidthMHz"`
	Channel         int     `json:"channel"`
}

// DeviceStatistics represents device statistics.
type DeviceStatistics struct {
	UptimeSec            int64   `json:"uptimeSec"`
	LastHeartbeatAt      string  `json:"lastHeartbeatAt"`
	NextHeartbeatAt      string  `json:"nextHeartbeatAt"`
	LoadAverage1Min      float64 `json:"loadAverage1Min"`
	LoadAverage5Min      float64 `json:"loadAverage5Min"`
	LoadAverage15Min     float64 `json:"loadAverage15Min"`
	CPUUtilizationPct    float64 `json:"cpuUtilizationPct"`
	MemoryUtilizationPct float64 `json:"memoryUtilizationPct"`
	Uplink               *struct {
		TxRateBps int64 `json:"txRateBps"`
		RxRateBps int64 `json:"rxRateBps"`
	} `json:"uplink,omitempty"`
	Interfaces *struct {
		Radios []struct {
			FrequencyGHz  float64 `json:"frequencyGHz"`
			TxRetriesPct  float64 `json:"txRetriesPct"`
		} `json:"radios,omitempty"`
	} `json:"interfaces,omitempty"`
}

// NetworkClient represents a network client.
type NetworkClient struct {
	Type           string `json:"type"` // WIRED or WIRELESS
	ID             string `json:"id"`
	Name           string `json:"name"`
	ConnectedAt    string `json:"connectedAt"`
	IPAddress      string `json:"ipAddress"`
	MACAddress     string `json:"macAddress"`
	UplinkDeviceID string `json:"uplinkDeviceId"`
	Access         *struct {
		Type string `json:"type"`
	} `json:"access,omitempty"`
}

// Network represents a UniFi network.
type Network struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Management string `json:"management"`
	Enabled    bool   `json:"enabled"`
	VlanID     int    `json:"vlanId"`
	ZoneID     string `json:"zoneId"`
	Default    bool   `json:"default"`
	Metadata   *struct {
		Origin       string `json:"origin"`
		Configurable bool   `json:"configurable"`
	} `json:"metadata,omitempty"`
}

// NetworkCreateRequest represents a request to create a network.
type NetworkCreateRequest struct {
	Name     string `json:"name"`
	Purpose  string `json:"purpose,omitempty"`
	VlanID   int    `json:"vlanId,omitempty"`
	Subnet   string `json:"subnet,omitempty"`
	DHCPMode string `json:"dhcpMode,omitempty"`
}

// NetworkReference represents a dependency reference for a network.
type NetworkReference struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// FirewallPolicy represents a firewall policy.
type FirewallPolicy struct {
	Enabled        bool   `json:"enabled"`
	Name           string `json:"name"`
	Index          int    `json:"index"`
	LoggingEnabled bool   `json:"loggingEnabled"`
	Action         struct {
		Type string `json:"type"`
	} `json:"action"`
	Source *struct {
		ZoneID        string `json:"zoneId"`
		TrafficFilter *struct {
			Type string `json:"type"`
		} `json:"trafficFilter,omitempty"`
	} `json:"source,omitempty"`
	Destination *struct {
		ZoneID        string `json:"zoneId"`
		TrafficFilter *struct {
			Type string `json:"type"`
		} `json:"trafficFilter,omitempty"`
	} `json:"destination,omitempty"`
	IPProtocolScope *struct {
		IPVersion string `json:"ipVersion"`
	} `json:"ipProtocolScope,omitempty"`
}

// FirewallPolicyCreateRequest represents a request to create a firewall policy.
type FirewallPolicyCreateRequest struct {
	Name              string `json:"name"`
	Enabled           bool   `json:"enabled"`
	Action            string `json:"action"`
	SourceZoneID      string `json:"sourceZoneId,omitempty"`
	DestinationZoneID string `json:"destinationZoneId,omitempty"`
	Protocol          string `json:"protocol,omitempty"`
	Description       string `json:"description,omitempty"`
}

// FirewallZone represents a firewall zone.
type FirewallZone struct {
	ID         string   `json:"id"`
	Name       string   `json:"name"`
	NetworkIDs []string `json:"networkIds"`
	Metadata   *struct {
		Origin       string `json:"origin"`
		Configurable bool   `json:"configurable"`
	} `json:"metadata,omitempty"`
}

// FirewallZoneCreateRequest represents a request to create a firewall zone.
type FirewallZoneCreateRequest struct {
	Name       string   `json:"name"`
	NetworkIDs []string `json:"networkIds,omitempty"`
}

// OrderingRequest represents a reorder request.
type OrderingRequest struct {
	IDs []string `json:"ids"`
}

// WiFiBroadcast represents a WiFi broadcast/SSID.
type WiFiBroadcast struct {
	Type    string `json:"type"`
	ID      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	Network *struct {
		Type string `json:"type"`
	} `json:"network,omitempty"`
	SecurityConfiguration *struct {
		Type string `json:"type"`
	} `json:"securityConfiguration,omitempty"`
	BroadcastingFrequenciesGHz []float64 `json:"broadcastingFrequenciesGHz"`
	BroadcastingDeviceFilter   *struct {
		Type      string   `json:"type"`
		DeviceIDs []string `json:"deviceIds"`
	} `json:"broadcastingDeviceFilter,omitempty"`
	Metadata *struct {
		Origin string `json:"origin"`
	} `json:"metadata,omitempty"`
}

// WiFiCreateRequest represents a request to create a WiFi broadcast.
type WiFiCreateRequest struct {
	Name      string `json:"name"`
	SSID      string `json:"ssid,omitempty"`
	Enabled   bool   `json:"enabled"`
	Security  string `json:"security,omitempty"`
	Password  string `json:"password,omitempty"`
	NetworkID string `json:"networkId,omitempty"`
	Band      string `json:"band,omitempty"`
	IsGuest   bool   `json:"isGuest,omitempty"`
}

// ACLRule represents an ACL rule.
type ACLRule struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Action      string `json:"action"`
	SourceMAC   string `json:"sourceMac"`
	Description string `json:"description"`
	Index       int    `json:"index"`
}

// ACLRuleCreateRequest represents a request to create an ACL rule.
type ACLRuleCreateRequest struct {
	Name        string `json:"name"`
	Enabled     bool   `json:"enabled"`
	Action      string `json:"action"`
	SourceMAC   string `json:"sourceMac,omitempty"`
	Description string `json:"description,omitempty"`
}

// DNSRecord represents a DNS record (the integration API returns DNS as records, not policies).
type DNSRecord struct {
	Type        string `json:"type"` // A_RECORD, etc.
	ID          string `json:"id"`
	Enabled     bool   `json:"enabled"`
	Domain      string `json:"domain"`
	IPv4Address string `json:"ipv4Address"`
	TTLSeconds  int    `json:"ttlSeconds"`
	Metadata    *struct {
		Origin string `json:"origin"`
	} `json:"metadata,omitempty"`
}

// DNSRecordCreateRequest represents a request to create a DNS record.
type DNSRecordCreateRequest struct {
	Type        string `json:"type,omitempty"`
	Enabled     bool   `json:"enabled"`
	Domain      string `json:"domain"`
	IPv4Address string `json:"ipv4Address,omitempty"`
	TTLSeconds  int    `json:"ttlSeconds,omitempty"`
}

// DeviceAction represents an action to perform on a device.
type DeviceAction struct {
	Action string `json:"action"`
}

// ClientAction represents an action to perform on a client.
type ClientAction struct {
	Action string `json:"action"`
}

// PortAction represents an action to perform on a port.
type PortAction struct {
	Action string `json:"action"`
}

// AdoptRequest represents a device adoption request.
type AdoptRequest struct {
	MAC string `json:"mac"`
}
