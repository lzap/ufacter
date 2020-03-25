package disk

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/lzap/ufacter/facts/common"
	"github.com/lzap/ufacter/lib/ufacter"
	d "github.com/shirou/gopsutil/disk"
)

var (
	reDevBlacklist = regexp.MustCompile("^(dm-[0-9]+|loop[0-9]+)$")
)

// getBlockDevices returns list of block devices
func getBlockDevices(all bool) ([]string, error) {
	blockDevs := []string{}
	targetDir := fmt.Sprintf("%v/block", common.GetHostSys())
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
	modelFilename := fmt.Sprintf("%s/block/%s/device/model",
		common.GetHostSys(), blockDevice)
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
	sizeFilename := fmt.Sprintf("%s/block/%s/size", common.GetHostSys(),
		blockDevice)
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
	vendorFilename := fmt.Sprintf("%s/block/%s/device/vendor",
		common.GetHostSys(), blockDevice)
	vendor, err := ioutil.ReadFile(vendorFilename)
	if err != nil {
		return "", err
	}
	vendor = bytes.TrimSuffix(vendor, []byte("\n"))
	vendor = bytes.TrimRight(vendor, " ")
	return fmt.Sprintf("%s", vendor), nil
}

func reportHumanReadable(facts chan<- ufacter.Fact, volatile bool, value uint64, mountpoint string, raw_key string, human_key string) error {
	human, unit, err := common.ConvertBytes(value)
	if err != nil {
		return err
	}
	facts <- ufacter.NewFact(value, false, "mountpoints", mountpoint, raw_key)
	facts <- ufacter.NewFact(fmt.Sprintf("%.2f %v", human, unit), false, "mountpoints", mountpoint, human_key)
	return nil
}

// ReportFacts returns related to HDDs
func ReportFacts(facts chan<- ufacter.Fact) error {
	partitions, err := d.Partitions(false)
	if err != nil {
		return err
	}

	for _, part := range partitions {
		usage, err := d.Usage(part.Mountpoint)
		if err != nil {
			return err
		}
		facts <- ufacter.NewStableFact(part.Device, "mountpoints", part.Mountpoint, "device")
		facts <- ufacter.NewStableFact(part.Fstype, "mountpoints", part.Mountpoint, "filesystem")
		facts <- ufacter.NewStableFact(strings.Split(part.Opts, ","), "mountpoints", part.Mountpoint, "options")
		facts <- ufacter.NewVolatileFact(fmt.Sprintf("%.2f%%", usage.UsedPercent), "mountpoints", part.Mountpoint, "capacity")
		reportHumanReadable(facts, false, usage.Total, part.Mountpoint, "size_bytes", "size")
		reportHumanReadable(facts, true, usage.Free, part.Mountpoint, "available_bytes", "available")
		reportHumanReadable(facts, true, usage.Used, part.Mountpoint, "used_bytes", "used")
	}

	blockDevs, err := getBlockDevices(false)
	if err != nil {
		return err
	}

	sort.Strings(blockDevs)
	for _, blockDevice := range blockDevs {
		size, err := getBlockDeviceSize(blockDevice)
		if err == nil {
			reportHumanReadable(facts, false, uint64(size), blockDevice, "size_bytes", "size")
		}

		model, err := getBlockDeviceModel(blockDevice)
		if err == nil {
			facts <- ufacter.NewStableFact(model, "disks", blockDevice, "model")
		}

		vendor, err := getBlockDeviceVendor(blockDevice)
		if err == nil {
			facts <- ufacter.NewStableFact(vendor, "disks", blockDevice, "vendor")
		}
	}

	facts <- ufacter.NewLastFact()
	return nil
}
