package bare_metal

type NetworkSize struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CIDR         string `json:"cidr"`
	NetworkType  string `json:"networkType"`
	AvailableIPs int    `json:"availableIps"`
	Price        Price  `json:"price"`
}

func (n NetworkSize) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            n.ID,
		"name":          n.Name,
		"cidr":          n.CIDR,
		"networkType":   n.NetworkType,
		"available_ips": n.AvailableIPs,
		"price":         n.Price.String(),
	}
}
