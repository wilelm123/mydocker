package network

import (
	"github.com/vishvananda/netlink"
	"net"
)

var (
	defaultNetworkPath = "/var/run/mydocker/network/network"
	drivers            = map[string]NetworkDriver{}
	networks           = map[string]*Network{}
)

type Endpoint struct {
	ID          string           `json:"id"`
	Device      netlink.Veth     `json:"dev"`
	IPAddress   net.IP           `json:"ip"`
	MacAddress  net.HardwareAddr `json:"mac"`
	Network     *Network
	PortMapping []string
}

type Network struct {
	Name    string
	IPRange *net.IPNet
	Driver  string
}

type NetworkDriver interface {
	Name() string
	Create(subnet string, name string) (*Network, error)
	Delete(network Network) error
	Connect(network *Network, endpoint *Endpoint) error
	Disconnect(network Network, endpoint *Endpoint) error
}

func (nw *Network) dump(dumpPath string) error {

}