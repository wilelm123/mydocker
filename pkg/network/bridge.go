package network

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"net"
	"os/exec"
	"strings"
	"time"
)

type BridgeNetworkDriver struct {
}

func (b *BridgeNetworkDriver) Name() string {
	return "bridge"
}

func (b *BridgeNetworkDriver) Create(subnet string, name string) (*Network, error) {
	ip, ipRange, _ := net.ParseCIDR(subnet)
	ipRange.IP = ip

	n := &Network{
		Name:    name,
		IPRange: ipRange,
		Driver:  b.Name(),
	}

	err := b.initBridge(n)
	if err != nil {
		log.Errorf("error init bridge: %v", err)
	}
	return n, err
}

func (b *BridgeNetworkDriver) Delete(network Network) error {
	bridgeName := network.Name
	brn, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}
	return netlink.LinkDel(brn)
}

func (b *BridgeNetworkDriver) Connect(network *Network, endpoint *Endpoint) error {
	bridgeName := network.Name
	brn, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return err
	}

	linkAttr := netlink.NewLinkAttrs()
	linkAttr.Name = endpoint.ID[:5]
	linkAttr.MasterIndex = brn.Attrs().Index

	endpoint.Device = netlink.Veth{
		LinkAttrs: linkAttr,
		PeerName:  "cif-" + endpoint.ID[:5],
	}

	if err = netlink.LinkAdd(&endpoint.Device); err != nil {
		return fmt.Errorf("error add endpoint device: %v", err)
	}

	if err = netlink.LinkSetUp(&endpoint.Device); err != nil {
		return fmt.Errorf("error add endpoint device: %v", err)
	}
	return nil
}

func (b *BridgeNetworkDriver) Disconnect(network Network, endpoint *Endpoint) error {
	return nil
}

func (b *BridgeNetworkDriver) initBridge(n *Network) error {
	bridgeName := n.Name
	if err := createBridgeInterface(bridgeName); err != nil {
		return fmt.Errorf("error add bridge: %s, error: %v", bridgeName, err)
	}

	gatewayIP := *n.IPRange
	gatewayIP.IP = n.IPRange.IP

	if err := setInterfaceIP(bridgeName, gatewayIP.String()); err != nil {
		return fmt.Errorf("error assigning address: %s on bridge: %s with an error: %v", gatewayIP, bridgeName, err)
	}

	if err := setInterfaceUP(bridgeName); err != nil {
		return fmt.Errorf("error set bridge up: %s, error: %v", bridgeName, err)
	}

	if err := setupIPTables(bridgeName, n.IPRange); err != nil {
		return fmt.Errorf("error setting iptables for %s, error: %v", bridgeName, err)
	}
	return nil
}

func (b *BridgeNetworkDriver) deleteBridge(n *Network) error {
	bridgeName := n.Name

	link, err := netlink.LinkByName(bridgeName)
	if err != nil {
		return fmt.Errorf("getting link with name %s failed: %v", bridgeName, err)
	}

	if err := netlink.LinkDel(link); err != nil {
		return fmt.Errorf("failed to remove bridge interface %s, error: %v", bridgeName, err)
	}
	return nil
}

func createBridgeInterface(bridgeName string) error {
	_, err := net.InterfaceByName(bridgeName)
	if err == nil || !strings.Contains(err.Error(), "no such network interface") {
		return err
	}
	linkAttr := netlink.NewLinkAttrs()
	linkAttr.Name = bridgeName

	b := &netlink.Bridge{LinkAttrs: linkAttr}
	if err := netlink.LinkAdd(b); err != nil {
		return fmt.Errorf("bridge creation failed for bridge %s: %v", bridgeName, err)
	}
	return nil
}

func setInterfaceUP(interfaceName string) error {
	iface, err := netlink.LinkByName(interfaceName)
	if err != nil {
		return fmt.Errorf("error retrieving a link named [ %s ]: %v", iface.Attrs().Name, err)
	}

	if err := netlink.LinkSetUp(iface); err != nil {
		return fmt.Errorf("error enabling interce for %s: %v", interfaceName, err)
	}
	return nil
}

func setInterfaceIP(name string, rawIP string) error {
	retries := 2
	var iface netlink.Link
	var err error
	for i := 0; i < retries; i++ {
		iface, err = netlink.LinkByName(name)
		if err == nil {
			break
		}
		log.Debugf("error retrieving new bridge netlink link [ %s ]... retrying", name)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		return fmt.Errorf("abandoning retrieving the new bridge link from netlink, Run [ ip link ] to troubleshooting the error: %v", err)
	}
	ipNet, err := netlink.ParseIPNet(rawIP)
	if err != nil {
		return err
	}

	addr := &netlink.Addr{IPNet: ipNet, Peer: ipNet, Label: "", Flags: 0, Scope: 0, Broadcast: nil}
	return netlink.AddrAdd(iface, addr)
}

func setupIPTables(bridgeName string, subnet *net.IPNet) error {
	iptablesCmd := fmt.Sprintf("-t nat -A ROUTING -s %s ! -o %s -j MASQUERADE", subnet.String(), bridgeName)
	cmd := exec.Command("iptables", strings.Split(iptablesCmd, " ")...)
	output, err := cmd.Output()
	if err != nil {
		log.Errorf("iptables Output: %s, error: %v", output, err)
	}
	return err
}
