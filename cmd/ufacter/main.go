package main

import (
	"flag"
	"io/ioutil"
	"os"
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
	"gopkg.in/yaml.v2"
)

func main() {
	conf := ufacter.Config{}
	modules := flag.String("modules", "cpu,mem,host,disk,net,route,link,ufacter", "Modules to run")
	yamlFormat := flag.Bool("yaml", false, "Print facts in YAML format")
	jsonFormat := flag.Bool("json", false, "Print facts in JSON format")
	noVolatile := flag.Bool("no-volatile", false, "Avoid facts that change often (e.g. free memory)")
	noExtended := flag.Bool("no-extended", false, "Avoid facts not found in the original facter")
	customFacts := flag.String("custom-facts", "", "Custom facts stored as YAML file")
	flag.Parse()

	if *yamlFormat == true {
		conf.Formatter = ufacter.NewYAMLFormatter()
	} else if *jsonFormat == true {
		conf.Formatter = ufacter.NewJSONFormatter()
	} else {
		// YAML is the default output in ufacter
		conf.Formatter = ufacter.NewYAMLFormatter()
	}

	// load custom facts first
	if _, err := os.Stat(*customFacts); err == nil {
		yamlMap := make(map[string]interface{})

		yamlString, err := ioutil.ReadFile(*customFacts)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(yamlString, &yamlMap)
		if err != nil {
			panic(err)
		}

		for key, value := range yamlMap {
			conf.Formatter.Add(ufacter.NewStableFact(value, key))
		}
	}

	// channel buffer hasn't measurable effect only for light formatters
	factsCh := make(chan ufacter.Fact, 1024)

	// slice of reporters (put your new reporter HERE)
	var reporters []func(facts chan<- ufacter.Fact, volatile bool, extended bool)
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
			go r(factsCh, !*noVolatile, !*noExtended)
		}

		// collect and wait for facts
		for f := range factsCh {
			if f.Name == nil {
				toClose--
			} else {
				if f.Value != nil && f.Value != "" {
					if (*noVolatile && f.Volatile) || (*noExtended && !f.Native) {
						// skip
					} else {
						conf.Formatter.Add(f)
					}
				}
			}
			if toClose <= 0 {
				break
			}
		}
		conf.Formatter.Finish()
	}
}
