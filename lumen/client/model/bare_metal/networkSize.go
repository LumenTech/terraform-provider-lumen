package bare_metal

type NetworkSize struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CIDR         string `json:"cidr"`
	AvailableIPs int    `json:"availableIps"`
	Price        Price  `json:"price"`
}

func (n NetworkSize) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":            n.ID,
		"name":          n.Name,
		"cidr":          n.CIDR,
		"available_ips": n.AvailableIPs,
		"price":         n.Price.String(),
	}
}
