package facter

import (
	dot "github.com/joeycumines/go-dotnotation/dotnotation"
	"github.com/lzap/ufacter/lib/formatter"
)

// Facter struct holds Facter-related attributes
type Facter struct {
	facts     map[string]interface{}
	formatter Formatter
}

// Config struct serves to pass Facter configuration
type Config struct {
	Formatter Formatter
}

// Formatter interface
type Formatter interface {
	Print(map[string]interface{}) error
}

// New returns new instance of Facter
func New(userConf *Config) *Facter {
	var conf *Config
	if userConf != nil {
		conf = userConf
	} else {
		conf = &Config{
			Formatter: formatter.NewFormatter(),
		}
	}
	f := &Facter{
		facts:     make(map[string]interface{}),
		formatter: conf.Formatter,
	}
	return f
}

// Add a fact
func (f *Facter) Add(k string, v interface{}) {
	dot.Set(f.facts, k, v)
}

// Add a fact
func (f *Facter) AddNode(k string) {
	dot.Set(f.facts, k, make(map[string]interface{}))
}

// Get returns value of given fact, if it exists
func (f *Facter) Get(k string) (interface{}, error) {
	value, err := dot.Get(f.facts, k)
	//value, ok := f.facts[k]
	return value, err
}

// Print prints-out facts by calling formatter
func (f *Facter) Print() {
	f.formatter.Print(f.facts)
}
