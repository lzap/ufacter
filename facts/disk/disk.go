package disk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	c "github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	d "github.com/shirou/gopsutil/disk"
)

var (
	reDevBlacklist = regexp.MustCompile("^(dm-[0-9]+|loop[0-9]+)$")
)

// getBlockDevices returns list of block devices
func getBlockDevices(all bool) ([]string, error) {
	blockDevs := []string{}
	targetDir := fmt.Sprintf("%v/block", c.GetHostSys())
	contents, err := ioutil.ReadDir(targetDir)
	if err != nil {
		return blockDevs, err
	}
	for _, v := range contents {
		if all == false {
			if reDevBlacklist.MatchString(v.Name()) {
				continue
			}
		}
		blockDevs = append(blockDevs, v.Name())
	}
	return blockDevs, nil
}

// getBlockDeviceModel returns model of block device as reported by Linux
// kernel.
func getBlockDeviceModel(blockDevice string) (string, error) {
	modelFilename := fmt.Sprintf("%s/block/%s/device/model", c.GetHostSys(), blockDevice)
	if _, err := os.Stat(modelFilename); err != nil {
		return "", nil
	}

	model, err := ioutil.ReadFile(modelFilename)
	if err != nil {
		return "", err
	}
	model = bytes.TrimSuffix(model, []byte("\n"))
	model = bytes.TrimSpace(model)
	return fmt.Sprintf("%s", model), nil
}

// getBlockDeviceSize returns size of block device as reported by Linux kernel
// multiplied by 512.
func getBlockDeviceSize(blockDevice string) (int64, error) {
	sizeFilename := fmt.Sprintf("%s/block/%s/size", c.GetHostSys(), blockDevice)
	if _, err := os.Stat(sizeFilename); err != nil {
		return 0, nil
	}

	size, err := ioutil.ReadFile(sizeFilename)
	if err != nil {
		return 0, err
	}
	sizeInt, err := strconv.ParseInt(fmt.Sprintf("%s",
		bytes.TrimSuffix(size, []byte("\n"))), 10, 64)
	if err != nil {
		return 0, err
	}
	return sizeInt * 512, nil
}

// getBlockDeviceVendor returns vendor of block device as reported by Linux
// kernel.
func getBlockDeviceVendor(blockDevice string) (string, error) {
	vendorFilename := fmt.Sprintf("%s/block/%s/device/vendor", c.GetHostSys(), blockDevice)
	if _, err := os.Stat(vendorFilename); err != nil {
		return "", nil
	}

	vendor, err := ioutil.ReadFile(vendorFilename)
	if err != nil {
		return "", err
	}
	vendor = bytes.TrimSuffix(vendor, []byte("\n"))
	vendor = bytes.TrimRight(vendor, " ")
	return fmt.Sprintf("%s", vendor), nil
}

func reportHumanReadable(facts chan<- ufacter.Fact, volatile bool, value uint64, mountpoint string, raw_key string, human_key string) error {
	human, unit, err := c.ConvertBytes(value)
	if err != nil {
		return err
	}
	facts <- ufacter.NewFact(value, false, "mountpoints", mountpoint, raw_key)
	facts <- ufacter.NewFact(fmt.Sprintf("%.2f %v", human, unit), false, "mountpoints", mountpoint, human_key)
	return nil
}

// ReportFacts returns related to HDDs
func ReportFacts(facts chan<- ufacter.Fact, volatile bool, extended bool) {
	start := time.Now()
	defer ufacter.SendLastFact(facts)

	partitions, err := d.Partitions(false)
	if err == nil {
		for _, part := range partitions {
			usage, err := d.Usage(part.Mountpoint)
			if err == nil {
				facts <- ufacter.NewStableFact(part.Device, "partitions", part.Device, "device")
				facts <- ufacter.NewStableFact(part.Fstype, "partitions", part.Device, "filesystem")
				facts <- ufacter.NewStableFact(strings.Split(part.Opts, ","), "partitions", part.Device, "options")
				facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", usage.UsedPercent), "partitions", part.Device, "capacity")
				facts <- ufacter.NewStableFact(usage.Total, "partitions", part.Device, "size_bytes")
				facts <- ufacter.NewStableFact(c.ConvertBytesAsString(usage.Total), "partitions", part.Device, "size")
				facts <- ufacter.NewVolatileFact(usage.Free, "partitions", part.Device, "available_bytes")
				facts <- ufacter.NewVolatileFact(c.ConvertBytesAsString(usage.Free), "partitions", part.Device, "available")
				facts <- ufacter.NewVolatileFact(usage.Used, "partitions", part.Device, "used_bytes")
				facts <- ufacter.NewVolatileFact(c.ConvertBytesAsString(usage.Used), "partitions", part.Device, "used")
			} else {
				c.LogError(facts, err, "disk", "usage")
			}
		}
	} else {
		c.LogError(facts, err, "disk", "partitions")
	}

	var sizeTotal uint64
	blockDevs, err := getBlockDevices(false)
	if err == nil {
		for _, blockDevice := range blockDevs {
			size, err := getBlockDeviceSize(blockDevice)
			sizeTotal += uint64(size)
			if err == nil {
				facts <- ufacter.NewStableFact(size, "disks", blockDevice, "size_bytes")
				facts <- ufacter.NewStableFact(c.ConvertBytesAsString(uint64(size)), "disks", blockDevice, "size")
			} else {
				c.LogError(facts, err, "disk", "block device size")
			}

			model, err := getBlockDeviceModel(blockDevice)
			if err == nil {
				facts <- ufacter.NewStableFact(model, "disks", blockDevice, "model")
			} else {
				c.LogError(facts, err, "disk", "block device model")
			}

			vendor, err := getBlockDeviceVendor(blockDevice)
			if err == nil {
				facts <- ufacter.NewStableFact(vendor, "disks", blockDevice, "vendor")
			} else {
				c.LogError(facts, err, "disk", "block device vendor")
			}

			ioc, err := d.IOCounters(blockDevice)
			if err == nil {
				facts <- ufacter.NewStableFact(ioc[blockDevice].Label, "disks", blockDevice, "label")
				facts <- ufacter.NewStableFact(ioc[blockDevice].SerialNumber, "disks", blockDevice, "serial")
			} else {
				c.LogError(facts, err, "disk", "block device iocounters")
			}
		}
		facts <- ufacter.NewStableFactEx(sizeTotal, "disks", "total_size_bytes")
		facts <- ufacter.NewVolatileFact(c.ConvertBytesAsString(sizeTotal), "disks", "total_size")
	} else {
		c.LogError(facts, err, "disk", "block devices")
	}

	ufacter.SendVolatileFactEx(facts, time.Since(start), "ufacter", "stats", "disk")
}
