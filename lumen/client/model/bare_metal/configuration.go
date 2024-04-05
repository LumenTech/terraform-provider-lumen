package bare_metal

import "fmt"

type Configuration struct {
	Name         string `json:"name"`
	DisplayName  string `json:"displayName"`
	Cores        int    `json:"cores"`
	Memory       string `json:"memory"`
	Storage      string `json:"storage"`
	Disks        int    `json:"disks"`
	Nics         int    `json:"nics"`
	Processors   int    `json:"processors"`
	MachineCount int    `json:"machineCount"`
	Price        Price  `json:"price"`
	Tier         *int   `json:"tier"`
}

func (c Configuration) ToMap() map[string]interface{} {
	tier := "None Set"
	if c.Tier != nil {
		tier = fmt.Sprintf("%d", *c.Tier)
	}

	return map[string]interface{}{
		"name":          c.Name,
		"display_name":  c.DisplayName,
		"cores":         c.Cores,
		"memory":        c.Memory,
		"storage":       c.Storage,
		"disks":         c.Disks,
		"nics":          c.Nics,
		"processors":    c.Processors,
		"machine_count": c.MachineCount,
		"price":         c.Price.String(),
		"tier":          tier,
	}
}
