package bare_metal

type Network struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	AccountID      string           `json:"accountId"`
	ServiceID      string           `json:"serviceId"`
	Location       string           `json:"location"`
	LocationID     string           `json:"locationId"`
	IPBlock        string           `json:"ipBlock"`
	IPV6Block      string           `json:"ipv6Block"`
	Gateway        string           `json:"gateway"`
	AvailableIPs   int              `json:"availableIps"`
	TotalIPs       int              `json:"totalIPs"`
	Type           string           `json:"type"`
	Status         string           `json:"status"`
	VRF            string           `json:"vrf,omitempty"`
	VRFDescription string           `json:"vrfDescription,omitempty"`
	Prices         []ComponentPrice `json:"prices"`
	Created        string           `json:"created"`
	Updated        string           `json:"updated"`
}

type NetworkProvisionRequest struct {
	Name           string `json:"name"`
	LocationID     string `json:"locationId"`
	NetworkSizeID  string `json:"networkSizeId,omitempty"`
	NetworkType    string `json:"networkType"`
	VRF            string `json:"vrf,omitempty"`
	VRFDescription string `json:"vrfDescription,omitempty"`
}

type AddNetworkRequest struct {
	NetworkId         string `json:"networkId"`
	AssignIPV6Address bool   `json:"assignIpv6Address,omitempty"`
}

type NetworkUpdateRequest struct {
	Name string `json:"name"`
}
