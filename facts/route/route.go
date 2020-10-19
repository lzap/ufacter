package route

import (
	"net"
	"time"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	s "github.com/shirou/gopsutil/net"
	n "github.com/vishvananda/netlink"
)

// ReportFacts adds route information
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	start := time.Now()
	defer ufacter.SendLastFact(facts)

	var primaryIPv4MAC string
	var primaryIPv6MAC string

	// primary IPv4 interface
	defaultRoutes, err := n.RouteGet(net.ParseIP("1.0.0.0"))
	if err == nil && len(defaultRoutes) > 0 {
		primaryLink, err := n.LinkByIndex(defaultRoutes[0].LinkIndex)
		if err == nil {
			primaryIPv4MAC = primaryLink.Attrs().HardwareAddr.String()
			facts <- ufacter.NewStableFact(primaryLink.Attrs().Name, "networking", "primary")
			facts <- ufacter.NewStableFactEx(primaryLink.Attrs().HardwareAddr.String(), "networking", "primary_mac")
		} else {
			c.LogError(facts, err, "link", "link by index")
		}
	} else {
		c.LogError(facts, err, "link", "IPv4 route")
	}

	// primary IPv6 interface
	defaultRoutes, err = n.RouteGet(net.ParseIP("100::"))
	if err == nil && len(defaultRoutes) > 0 {
		primaryLink, err := n.LinkByIndex(defaultRoutes[0].LinkIndex)
		if err == nil {
			primaryIPv6MAC = primaryLink.Attrs().HardwareAddr.String()
			facts <- ufacter.NewStableFact(primaryLink.Attrs().Name, "networking", "primary6")
			facts <- ufacter.NewStableFactEx(primaryLink.Attrs().HardwareAddr.String(), "networking", "primary6_mac")
		} else {
			c.LogError(facts, err, "link", "link by index")
		}
	} else {
		c.LogError(facts, err, "link", "IPv6 route")
	}

	// networking.(mac,mtu,ip,ip6,netmask,netmask6,network,network6) + (mac6,mtu6) extended
	netIfaces, err := s.Interfaces()
	if err != nil {
		c.LogError(facts, err, "route", "interfaces")
		return
	}

	for _, v := range netIfaces {
		if v.HardwareAddr == primaryIPv4MAC && len(primaryIPv6MAC) > 0 {
			// IPv4
			facts <- ufacter.NewStableFact(v.HardwareAddr, "networking", "mac")
			facts <- ufacter.NewStableFact(v.MTU, "networking", "mtu")

			for _, ipAddr := range v.Addrs {
				parsedIP, parsedNet, err := net.ParseCIDR(ipAddr.Addr)
				if err != nil {
					c.LogError(facts, err, "route", "parse CIDR")
					continue
				}
				if ip4 := parsedIP.To4(); ip4 != nil {
					facts <- ufacter.NewStableFactEx(ipAddr.Addr, "networking", "cidr")
					facts <- ufacter.NewStableFact(parsedIP.String(), "networking", "ip")
					facts <- ufacter.NewStableFact(parsedNet.IP.String(), "networking", "network")
					facts <- ufacter.NewStableFact(c.IPMaskToString4(parsedNet.Mask), "networking", "netmask")
				}
			}
		} else if v.HardwareAddr == primaryIPv6MAC && len(primaryIPv6MAC) > 0 {
			// IPv6
			facts <- ufacter.NewStableFactEx(v.HardwareAddr, "networking", "mac6")
			facts <- ufacter.NewStableFactEx(v.MTU, "networking", "mtu6")

			for _, ipAddr := range v.Addrs {
				parsedIP, parsedNet, err := net.ParseCIDR(ipAddr.Addr)
				if err != nil {
					c.LogError(facts, err, "route", "parse CIDR")
					continue
				}
				if ip4 := parsedIP.To4(); ip4 == nil {
					facts <- ufacter.NewStableFactEx(ipAddr.Addr, "networking", "cidr6")
					facts <- ufacter.NewStableFact(parsedIP.String(), "networking", "ip6")
					facts <- ufacter.NewStableFact(parsedNet.IP.String(), "networking", "network6")
					facts <- ufacter.NewStableFact(c.IPMaskToString4(parsedNet.Mask), "networking", "netmask6")
				}
			}
		}
	}

	ufacter.SendVolatileFactEx(facts, time.Since(start), "ufacter", "stats", "route")
}
