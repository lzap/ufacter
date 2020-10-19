package net

import (
	"net"
	"regexp"
	"strings"
	"time"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	n "github.com/shirou/gopsutil/net"
)

var (
	reIPv4 = regexp.MustCompile("^[0-9]+\\.")
)

type stringMap map[string]string

// ReportFacts gathers network related facts
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	start := time.Now()
	defer ufacter.SendLastFact(facts)

	netIfaces, err := n.Interfaces()
	if err != nil {
		c.LogError(facts, err, "net", "interfaces")
		return
	}

	var ifaces []string
	for _, v := range netIfaces {
		ifName := strings.ToLower(v.Name)
		ifaces = append(ifaces, ifName)
		if v.HardwareAddr != "" {
			facts <- ufacter.NewStableFact(v.HardwareAddr, "networking", "interfaces", ifName, "mac")
		}
		facts <- ufacter.NewStableFact(v.MTU, "networking", "interfaces", ifName, "mtu")

		bindings := make([]stringMap, 0)
		bindings6 := make([]stringMap, 0)
		for _, ipAddr := range v.Addrs {
			parsedIP, parsedNet, err := net.ParseCIDR(ipAddr.Addr)
			if err != nil {
				c.LogError(facts, err, "net", "parse CIDR")
				continue
			}
			b := make(stringMap)
			if extended {
				b["cidr"] = ipAddr.Addr
			}
			b["address"] = parsedIP.String()
			b["network"] = parsedNet.IP.String()
			if ip4 := parsedIP.To4(); ip4 != nil {
				b["netmask"] = c.IPMaskToString4(parsedNet.Mask)
				bindings = append(bindings, b)
			} else {
				b["netmask"] = c.IPMaskToString6(parsedNet.Mask)
				bindings6 = append(bindings6, b)
			}
			if len(bindings) > 0 {
				facts <- ufacter.NewStableFact(bindings, "networking", "interfaces", ifName, "bindings")
			}
			if len(bindings6) > 0 {
				facts <- ufacter.NewStableFact(bindings6, "networking", "interfaces", ifName, "bindings6")
			}
		}
	}

	ufacter.SendVolatileFactEx(facts, time.Since(start), "ufacter", "stats", "net")
}
