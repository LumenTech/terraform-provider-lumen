package bare_metal

type Network struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	AccountID    string           `json:"accountId"`
	ServiceID    string           `json:"serviceId"`
	Location     string           `json:"location"`
	LocationID   string           `json:"locationId"`
	IPBlock      string           `json:"ipBlock"`
	Gateway      string           `json:"gateway"`
	AvailableIPs int              `json:"availableIps"`
	TotalIPs     int              `json:"totalIPs"`
	Type         string           `json:"type"`
	Status       string           `json:"status"`
	Prices       []ComponentPrice `json:"prices"`
	Created      string           `json:"created"`
	Updated      string           `json:"updated"`
}

type NetworkUpdateRequest struct {
	NetworkID string `json:"networkId"`
	Name      string `json:"name"`
	AccountID string `json:"accountId"`
}
