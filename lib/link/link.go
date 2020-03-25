package link

import (
	"fmt"
	"regexp"

	n "github.com/vishvananda/netlink"
)

var (
	reIPv4 = regexp.MustCompile("^[0-9]+\\.")
)

// exported
type Facter interface {
	Add(string, interface{})
}

func idToName(id int) string {
	if id == 0 {
		return ""
	}

	link, err := n.LinkByIndex(id)
	if err != nil {
		return ""
	}

	return link.Attrs().Name
}

// exported
func GetLinkFacts(f Facter) error {
	links, err := n.LinkList()
	if err != nil {
		return err
	}

	for _, link := range links {
		device := link.Attrs().Name

		f.Add(fmt.Sprintf("link.%s.type", device), link.Type())
		if len(link.Attrs().HardwareAddr.String()) > 0 {
			f.Add(fmt.Sprintf("link.%s.mac", device), link.Attrs().HardwareAddr)
		}
		if link.Attrs().ParentIndex != 0 {
			f.Add(fmt.Sprintf("link.%s.parent", device), idToName(link.Attrs().ParentIndex))
		}
		if link.Attrs().MasterIndex != 0 {
			f.Add(fmt.Sprintf("link.%s.master", device), idToName(link.Attrs().MasterIndex))
		}
		if link.Attrs().Slave != nil {
			f.Add(fmt.Sprintf("link.%s.slave", device), link.Attrs().Slave.SlaveType())
		}
		if link.Type() == "vlan" {
			vlan := link.(*n.Vlan)
			f.Add(fmt.Sprintf("link.%s.vlan.id", device), vlan.VlanId)
			f.Add(fmt.Sprintf("link.%s.vlan.protocol", device), vlan.VlanProtocol)
		}
		if link.Type() == "vxlan" {
			vxlan := link.(*n.Vxlan)
			f.Add(fmt.Sprintf("link.%s.vxlan.id", device), vxlan.VxlanId)
		}
		if link.Type() == "bond" {
			bond := link.(*n.Bond)
			f.Add(fmt.Sprintf("link.%s.bond.mode", device), bond.Mode)
		}
		if link.Type() == "veth" {
			veth := link.(*n.Veth)
			f.Add(fmt.Sprintf("link.%s.peer.name", device), veth.PeerName)
			f.Add(fmt.Sprintf("link.%s.peer.mac", device), veth.PeerHardwareAddr)
		}
	}
	return nil
}
