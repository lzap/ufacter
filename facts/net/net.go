package net

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/lzap/ufacter/lib/ufacter"
	n "github.com/shirou/gopsutil/net"
)

var (
	reIPv4 = regexp.MustCompile("^[0-9]+\\.")
)

type stringMap map[string]string

// Gathers network related facts
func ReportFacts(facts chan<- ufacter.Fact) error {
	netIfaces, err := n.Interfaces()
	if err != nil {
		return err
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
			parsed_ip, parsed_net, err := net.ParseCIDR(ipAddr.Addr)
			if err != nil {
				return err
			}
			b := make(stringMap)
			b["cidr"] = ipAddr.Addr
			b["address"] = parsed_ip.String()
			b["network"] = parsed_net.IP.String()
			var maskBuilder strings.Builder
			// TODO https://github.com/golang/go/issues/38097
			if !strings.ContainsRune(ipAddr.Addr, ':') {
				// IPv4
				maskBuilder.Grow(16)
				for i, b := range parsed_net.Mask {
					fmt.Fprintf(&maskBuilder, "%d", b)
					if i < 3 {
						maskBuilder.WriteString(".")
					}
				}
				b["netmask"] = maskBuilder.String()
				bindings = append(bindings, b)
			} else {
				// IPv6
				maskBuilder.Grow(41)
				for i, b := range parsed_net.Mask {
					fmt.Fprintf(&maskBuilder, "%x", b)
					if i%2 == 1 && i < 15 {
						maskBuilder.WriteString(":")
					}
				}
				b["netmask"] = maskBuilder.String()
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

	facts <- ufacter.NewLastFact()
	return nil
}
