package route

import (
	"net"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	n "github.com/vishvananda/netlink"
)

// ReportFacts adds route information
func ReportFacts(facts chan<- ufacter.Fact) error {
	var primaryIPv4Device string
	var primaryIPv6Device string

	// primary IPv4 interface
	defaultRoutes, err := n.RouteGet(net.ParseIP("1.0.0.0"))
	if err == nil && len(defaultRoutes) > 0 {
		primaryLink, err := n.LinkByIndex(defaultRoutes[0].LinkIndex)
		if err != nil {
			return err
		}
		primaryIPv4Device = primaryLink.Attrs().Name
		facts <- ufacter.NewStableFact(primaryIPv4Device, "network", "primary")
	} else {
		c.LogError(facts, "link: IPv6 route", err)
	}

	// primary IPv6 interface
	defaultRoutes, err = n.RouteGet(net.ParseIP("100::"))
	if err == nil && len(defaultRoutes) > 0 {
		primaryLink, err := n.LinkByIndex(defaultRoutes[0].LinkIndex)
		if err != nil {
			return err
		}
		primaryIPv6Device = primaryLink.Attrs().Name
		facts <- ufacter.NewStableFact(primaryIPv6Device, "network", "primary6")
	} else {
		c.LogError(facts, "link: IPv6 route", err)
	}

	// TODO: adding routing tables (would that be useful?)

	facts <- ufacter.NewLastFact()
	return nil
}
