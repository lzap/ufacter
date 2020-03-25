package host

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/lzap/ufacter/lib/ufacter"
	h "github.com/shirou/gopsutil/host"
)

// capitalize the first letter of given string
func capitalize(label string) string {
	firstLetter := strings.SplitN(label, "", 2)
	if len(firstLetter) < 1 {
		return label
	}
	return fmt.Sprintf("%v%v", strings.ToUpper(firstLetter[0]),
		strings.TrimPrefix(label, firstLetter[0]))
}

// int8ToString converts [65]int8 in syscall.Utsname to string
func int8ToString(bs [65]int8) string {
	b := make([]byte, len(bs))
	for i, v := range bs {
		if v < 0 {
			b[i] = byte(256 + int(v))
		} else {
			b[i] = byte(v)
		}
	}
	return strings.TrimRight(string(b), "\x00")
}

// GetHostFacts gathers facts related to Host
func ReportFacts(facts chan<- ufacter.Fact) error {
	hostInfo, err := h.Info()
	if err != nil {
		return err
	}

	facts <- ufacter.NewStableFact(hostInfo.Hostname, "networking", "fqdn")
	splitted := strings.SplitN(hostInfo.Hostname, ".", 2)
	var hostname *string
	if len(splitted) > 1 {
		hostname = &splitted[0]
		facts <- ufacter.NewStableFact(splitted[1], "networking", "domain")
	} else {
		hostname = &hostInfo.Hostname
	}
	facts <- ufacter.NewStableFact(*hostname, "networking", "hostname")

	var isVirtual bool
	if hostInfo.VirtualizationRole == "host" {
		isVirtual = false
	} else {
		isVirtual = true
	}
	facts <- ufacter.NewStableFact(isVirtual, "is_virtual")
	if hostInfo.VirtualizationRole == "host" {
		facts <- ufacter.NewStableFact("physical", "virtual")
	} else {
		facts <- ufacter.NewStableFact(hostInfo.VirtualizationSystem, "virtual")
	}

	facts <- ufacter.NewStableFact(capitalize(hostInfo.OS), "kernel")
	var uname syscall.Utsname
	err = syscall.Uname(&uname)
	if err == nil {
		kernelRelease := int8ToString(uname.Release)
		kernelVersion := strings.Split(kernelRelease, "-")[0]
		kvSplitted := strings.Split(kernelVersion, ".")
		facts <- ufacter.NewStableFact(kernelRelease, "kernelrelease")
		facts <- ufacter.NewStableFact(kernelVersion, "kernelversion")
		facts <- ufacter.NewStableFact(strings.Join(kvSplitted[0:2], "."), "kernelmajversion")
	}

	facts <- ufacter.NewStableFact(hostInfo.KernelArch, "os", "architecture")
	facts <- ufacter.NewStableFact(hostInfo.PlatformFamily, "os", "family")
	facts <- ufacter.NewStableFact(hostInfo.KernelArch, "os", "hardware")
	facts <- ufacter.NewStableFact(hostInfo.Platform, "os", "name")
	facts <- ufacter.NewStableFact(hostInfo.PlatformVersion, "os", "release", "full")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "release", "major")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "release", "minor")

	// LSB parsing: https://github.com/shirou/gopsutil/blob/master/host/host_linux.go
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "codename")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "description")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "id")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "specification")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "release", "full")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "release", "major")
	//facts <- ufacter.NewStableFact(hostInfo.OS, "os", "distro", "release", "minor")

	facts <- ufacter.NewVolatileFact(hostInfo.Uptime, "system_uptime", "seconds")
	facts <- ufacter.NewVolatileFact(hostInfo.Uptime/60/60, "system_uptime", "hours")
	facts <- ufacter.NewVolatileFact(hostInfo.Uptime/60/60/24, "system_uptime", "days")
	facts <- ufacter.NewVolatileFact(fmt.Sprintf("%d days", hostInfo.Uptime/60/60/24), "system_uptime", "uptime")
	facts <- ufacter.NewStableFactEx(hostInfo.BootTime, "system_uptime", "boot_time")

	envPath := os.Getenv("PATH")
	if envPath != "" {
		facts <- ufacter.NewStableFact(envPath, "path")
	}

	tz, _ := time.Now().Zone()
	facts <- ufacter.NewStableFact(tz, "timezone")

	facts <- ufacter.NewLastFact()
	return nil
}
