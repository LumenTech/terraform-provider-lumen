package bare_metal

type Server struct {
	ID            string              `json:"id"`
	Name          string              `json:"name"`
	MachineID     string              `json:"machineId"`
	MachineName   string              `json:"machineName"`
	LocationID    string              `json:"locationId"`
	Location      string              `json:"location"`
	Configuration ServerConfiguration `json:"configuration"`
	OSImage       string              `json:"osImage"`
	Networks      []ServerNetwork     `json:"networks"`
	Status        string              `json:"status"`
	StatusMessage string              `json:"statusMessage"`
	Disks         []Disk              `json:"disks"`
	BootDisk      string              `json:"bootDisk"`
	ServiceID     string              `json:"serviceId"`
	Prices        []ComponentPrice    `json:"prices"`
	AccountID     string              `json:"accountId"`
	Created       string              `json:"created"`
	Updated       string              `json:"updated"`
}

func (s Server) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   s.Name,
		"status": s.Status,
	}
}

type ServerConfiguration struct {
	Name       string `json:"name"`
	Cores      int    `json:"cores"`
	Memory     string `json:"memory"`
	Storage    string `json:"storage"`
	Disks      int    `json:"disks"`
	NICs       int    `json:"nics"`
	Processors int    `json:"processors"`
}

type ServerNetwork struct {
	ID            string `json:"id"`
	NetworkID     string `json:"networkId"`
	NetworkName   string `json:"networkName"`
	NetworkType   string `json:"networkType"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	IP            string `json:"ip"`
	VLAN          string `json:"vlan"`
}

type Disk struct {
	Boot     bool   `json:"boot"`
	DiskType string `json:"disk_type"`
	Path     string `json:"path"`
	Size     int64  `json:"size"`
}

type ComponentPrice struct {
	Type  string `json:"type"`
	Price Price  `json:"price"`
}
