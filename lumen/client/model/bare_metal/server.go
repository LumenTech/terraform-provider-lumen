package bare_metal

type Server struct {
	ID             string                   `json:"id"`
	Name           string                   `json:"name"`
	MachineID      string                   `json:"machineId"`
	MachineName    string                   `json:"machineName"`
	LocationID     string                   `json:"locationId"`
	Location       string                   `json:"location"`
	Configuration  ServerConfiguration      `json:"configuration"`
	OSImage        string                   `json:"osImage"`
	Networks       []ServerNetwork          `json:"networks"`
	Status         string                   `json:"status"`
	StatusMessage  string                   `json:"statusMessage"`
	Disks          []map[string]interface{} `json:"disks"`
	BootDisk       string                   `json:"bootDisk"`
	ServiceID      string                   `json:"serviceId"`
	Prices         []ComponentPrice         `json:"prices"`
	AccountID      string                   `json:"accountId"`
	Created        string                   `json:"created"`
	Updated        string                   `json:"updated"`
	Hyperthreading bool                     `json:"hyperthreading"`
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
	ID                string `json:"id"`
	NetworkID         string `json:"networkId"`
	NetworkName       string `json:"networkName"`
	NetworkType       string `json:"networkType"`
	Status            string `json:"status"`
	StatusMessage     string `json:"statusMessage"`
	IP                string `json:"ip"`
	IPV6              string `json:"ipv6"`
	VLAN              string `json:"vlan"`
	AssignIPV6Address bool   `json:"assignIpv6Address"`
}

type ComponentPrice struct {
	Type  string `json:"type"`
	Price Price  `json:"price"`
}

type ServerProvisionRequest struct {
	Name              string                   `json:"name"`
	LocationID        string                   `json:"locationId"`
	Configuration     string                   `json:"configuration"`
	OSImage           string                   `json:"osImage"`
	NetworkID         string                   `json:"networkId,omitempty"`
	NetworkRequest    *NetworkProvisionRequest `json:"networkRequest,omitempty"`
	Credentials       Credentials              `json:"credentials"`
	AssignIPV6Address bool                     `json:"assignIpv6Address,omitempty"`
	Hyperthreading    bool                     `json:"hyperthreading,omitempty"`
}

type Credentials struct {
	Username  string `json:"username"`
	Password  string `json:"password,omitempty"`
	PublicKey string `json:"publicKey,omitempty"`
}

type ServerUpdateRequest struct {
	Name string `json:"name"`
}
