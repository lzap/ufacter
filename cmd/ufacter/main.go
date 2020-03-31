package main

import (
	"flag"
	"strings"

	"github.com/lzap/ufacter/facts/cpu"
	"github.com/lzap/ufacter/facts/disk"
	"github.com/lzap/ufacter/facts/host"
	"github.com/lzap/ufacter/facts/link"
	"github.com/lzap/ufacter/facts/mem"
	"github.com/lzap/ufacter/facts/net"
	"github.com/lzap/ufacter/facts/route"
	fufacter "github.com/lzap/ufacter/facts/ufacter"
	"github.com/lzap/ufacter/lib/ufacter"
)

func main() {
	conf := ufacter.Config{}
	modules := flag.String("modules", "cpu,mem,host,disk,net,route,link,ufacter", "Modules to run")
	yamlFormat := flag.Bool("yaml", false, "Print facts in YAML format")
	jsonFormat := flag.Bool("json", false, "Print facts in JSON format")
	flag.Parse()

	if *yamlFormat == true {
		conf.Formatter = ufacter.NewYAMLFormatter()
	} else if *jsonFormat == true {
		conf.Formatter = ufacter.NewJSONFormatter()
	} else {
		// YAML is the default output in ufacter
		conf.Formatter = ufacter.NewYAMLFormatter()
	}

	// channel buffer hasn't measurable effect only for light formatters
	factsCh := make(chan ufacter.Fact, 1024)

	// slice of reporters (put your new reporter HERE)
	var reporters []func(facts chan<- ufacter.Fact)
	for _, mod := range strings.Split(*modules, ",") {
		switch mod {
		case "cpu":
			reporters = append(reporters, cpu.ReportFacts)
		case "mem":
			reporters = append(reporters, mem.ReportFacts)
		case "link":
			reporters = append(reporters, link.ReportFacts)
		case "route":
			reporters = append(reporters, route.ReportFacts)
		case "host":
			reporters = append(reporters, host.ReportFacts)
		case "net":
			reporters = append(reporters, net.ReportFacts)
		case "disk":
			reporters = append(reporters, disk.ReportFacts)
		case "ufacter":
			reporters = append(reporters, fufacter.ReportFacts)
		}
	}
	toClose := len(reporters)

	if toClose > 0 {
		// start all reporters
		for _, r := range reporters {
			go r(factsCh)
		}

		// collect and wait for facts
		for f := range factsCh {
			if f.Name == nil {
				toClose--
			} else {
				if f.Value != nil && f.Value != "" {
					conf.Formatter.Add(f)
				}
			}
			if toClose <= 0 {
				break
			}
		}
		conf.Formatter.Finish()
	}
}
