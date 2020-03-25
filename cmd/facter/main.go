package main

import (
	"flag"

	"github.com/zstyblik/go-facter/lib/cpu"
	"github.com/zstyblik/go-facter/lib/disk"
	"github.com/zstyblik/go-facter/lib/facter"
	"github.com/zstyblik/go-facter/lib/formatter"
	"github.com/zstyblik/go-facter/lib/host"
	"github.com/zstyblik/go-facter/lib/link"
	"github.com/zstyblik/go-facter/lib/mem"
	"github.com/zstyblik/go-facter/lib/net"
)

func main() {
	conf := facter.Config{}
	ptFormat := flag.Bool("plaintext", false,
		"Emit facts as key => value pairs")
	kvFormat := flag.Bool("keyvalue", false,
		"Emit facts as key:value pairs")
	jsonFormat := flag.Bool("json", false,
		"Emit facts as a JSON")
	flag.Parse()

	if *ptFormat == true {
		conf.Formatter = formatter.NewFormatter()
	} else if *kvFormat == true {
		conf.Formatter = formatter.NewKeyValueFormatter()
	} else if *jsonFormat == true {
		conf.Formatter = formatter.NewJSONFormatter()
	} else {
		conf.Formatter = formatter.NewFormatter()
	}

	facter := facter.New(&conf)
	_ = cpu.GetCPUFacts(facter)
	_ = disk.GetDiskFacts(facter)
	_ = host.GetHostFacts(facter)
	_ = mem.GetMemoryFacts(facter)
	_ = net.GetNetFacts(facter)
	_ = link.GetLinkFacts(facter)
	facter.Print()
}
