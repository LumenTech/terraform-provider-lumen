package bare_metal

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
	Tier         int    `json:"tier"`
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
		"tier":          c.Tier,
	}
}
