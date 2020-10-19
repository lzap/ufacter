package common

import (
	"fmt"
	"math"
	"net"
	"os"
	"strings"

	"github.com/lzap/ufacter/lib/ufacter"
)

var (
	// ByteUnits is a k=>v map of units for conversion
	ByteUnits = map[int]string{
		0: "bytes",
		1: "kB",
		2: "MiB",
		3: "GiB",
		4: "TiB",
	}
)

// LogError writes error into syslog (or logfile) and also into fact "ufacter.errors"
func LogError(facts chan<- ufacter.Fact, err error, keys ...string) {
	// TODO: journald or syslog
	facts <- ufacter.NewStableFact(err, append([]string{"ufacter", "errors"}, keys...)...)
}

// ConvertBytes converts bytes to the highest possible unit
func ConvertBytes(in uint64) (float64, string, error) {
	out := float64(in)
	idx := 0
	maxIdx := len(ByteUnits)
	for idx < maxIdx && out > 1 {
		tmp := float64(out) / 1024
		if tmp < 1 {
			break
		}
		out = tmp
		idx++
	}
	return out, ByteUnits[idx], nil
}

// ConvertBytesTo converts bytes to the specified unit
func ConvertBytesTo(in uint64, maxUnit string) (float64, string, error) {
	if maxUnit == "" {
		return 0, "", fmt.Errorf("Given maximum unit is invalid.")
	}
	out := float64(in)
	idx := 0
	maxIdx := len(ByteUnits)
	for idx < maxIdx && maxUnit != ByteUnits[idx] {
		out = float64(out) / 1024
		idx++
	}
	return out, ByteUnits[idx], nil
}

// ConvertNetmask converts CIDR (netmask) to Netmask
func ConvertNetmask(in uint8) (string, error) {
	if in > 32 {
		return "", fmt.Errorf("Invalid Netmask given.")
	}
	octets := map[uint8]uint8{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
	}
	var idx uint8 = 1
	for in > 0 && idx < 5 {
		if (in / 8) > 0 {
			in = in - 8
			octets[idx] = 255
		} else {
			mod := in % 8
			octets[idx] = 255 - uint8(math.Pow(2, float64(8-mod))) + 1
			in = 0
		}
		idx++
	}
	return fmt.Sprintf("%d.%d.%d.%d", octets[1], octets[2], octets[3],
		octets[4]), nil
}

// IPMaskToString4 converts IPMask to IPv4 string
func IPMaskToString4(mask net.IPMask) string {
	var maskBuilder strings.Builder
	maskBuilder.Grow(16)
	for i, b := range mask {
		fmt.Fprintf(&maskBuilder, "%d", b)
		if i < 3 {
			maskBuilder.WriteString(".")
		}
	}
	return net.ParseIP(maskBuilder.String()).String()
}

// IPMaskToString6 converts IPMask to IPv6 string
func IPMaskToString6(mask net.IPMask) string {
	var maskBuilder strings.Builder
	maskBuilder.Grow(41)
	for i, b := range mask {
		fmt.Fprintf(&maskBuilder, "%x", b)
		if i%2 == 1 && i < 15 {
			maskBuilder.WriteString(":")
		}
	}
	return net.ParseIP(maskBuilder.String()).String()
}

func GetHostEtc() string {
	host_etc := os.Getenv("HOST_ETC")
	if host_etc == "" {
		host_etc = "/etc"
	}
	return host_etc
}

func GetHostSys() string {
	host_sys := os.Getenv("HOST_SYS")
	if host_sys == "" {
		host_sys = "/sys"
	}
	return host_sys
}

func GetHostProc() string {
	host_proc := os.Getenv("HOST_PROC")
	if host_proc == "" {
		host_proc = "/proc"
	}
	return host_proc
}
