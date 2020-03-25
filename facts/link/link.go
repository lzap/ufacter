package link

import (
	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	n "github.com/vishvananda/netlink"
)

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

// ReportFacts adds link information
func ReportFacts(facts chan<- ufacter.Fact) error {
	links, err := n.LinkList()
	if err == nil {
		for _, link := range links {
			device := link.Attrs().Name

			facts <- ufacter.NewStableFact(link.Type(), "link", device, "type")
			if len(link.Attrs().HardwareAddr.String()) > 0 {
				facts <- ufacter.NewStableFact(link.Attrs().HardwareAddr.String(), "link", device, "mac")
			}
			if link.Attrs().ParentIndex != 0 {
				facts <- ufacter.NewStableFact(idToName(link.Attrs().ParentIndex), "link", device, "parent")
			}
			if link.Attrs().MasterIndex != 0 {
				facts <- ufacter.NewStableFact(idToName(link.Attrs().MasterIndex), "link", device, "master")
			}
			if link.Attrs().Slave != nil {
				facts <- ufacter.NewStableFact(link.Attrs().Slave.SlaveType(), "link", device, "slave")
			}
			if link.Type() == "vlan" {
				vlan := link.(*n.Vlan)
				facts <- ufacter.NewStableFact(vlan.VlanId, "link", device, "vlan", "id")
				facts <- ufacter.NewStableFact(vlan.VlanProtocol, "link", device, "vlan", "protocol")
			}
			if link.Type() == "vxlan" {
				vxlan := link.(*n.Vxlan)
				facts <- ufacter.NewStableFact(vxlan.VxlanId, "link", device, "vxlan", "id")
			}
			if link.Type() == "bond" {
				bond := link.(*n.Bond)
				facts <- ufacter.NewStableFact(bond.Mode, "link", device, "bond", "mode")
			}
			if link.Type() == "veth" {
				veth := link.(*n.Veth)
				facts <- ufacter.NewStableFact(veth.PeerName, "link", device, "peer", "name")
				facts <- ufacter.NewStableFact(veth.PeerHardwareAddr, "link", device, "peer", "mac")
			}
		}
	} else {
		c.LogError(facts, "link: getting list", err)
	}

	facts <- ufacter.NewLastFact()
	return nil
}
