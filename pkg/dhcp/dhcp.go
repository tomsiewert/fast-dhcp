package dhcp

type DHCPServer struct {
	Hostname  string `json:"hostname"`
	Interface string `json:"interface"`
	Port      int    `json:"port"`
}

func NewDHCPServer(hostname string, iface string, port int) *DHCPServer {
	return &DHCPServer{
		Hostname:  hostname,
		Interface: iface,
		Port:      port,
	}
}
