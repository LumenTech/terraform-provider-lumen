package bare_metal

type NetworkSizes []NetworkSize

func (n NetworkSizes) ToMapList() []map[string]interface{} {
	mapList := make([]map[string]interface{}, len(n))
	for idx, loc := range n {
		mapList[idx] = loc.ToMap()
	}
	return mapList
}

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
