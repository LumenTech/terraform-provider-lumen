package bare_metal

type Configurations []Configuration

func (configs Configurations) ToMapList() []map[string]interface{} {
	mapList := make([]map[string]interface{}, len(configs))
	for idx, loc := range configs {
		mapList[idx] = loc.ToMap()
	}
	return mapList
}

type Configuration struct {
	Name         string `json:"name"`
	Cores        int    `json:"cores"`
	Memory       string `json:"memory"`
	Storage      string `json:"storage"`
	Disks        int    `json:"disks"`
	Nics         int    `json:"nics"`
	Processors   int    `json:"processors"`
	MachineCount int    `json:"machineCount"`
	Price        Price  `json:"price"`
}

func (c Configuration) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":          c.Name,
		"cores":         c.Cores,
		"memory":        c.Memory,
		"storage":       c.Storage,
		"disks":         c.Disks,
		"nics":          c.Nics,
		"processors":    c.Processors,
		"machine_count": c.MachineCount,
		"price":         c.Price.String(),
	}
}
