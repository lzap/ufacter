package ufacter

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// YAMLFormatter prints-out facts in YAML format
type YAMLFormatter struct {
	data map[string]interface{}
}

// NewYAMLFormatter returns new YAML formatter
func NewYAMLFormatter() *YAMLFormatter {
	return &YAMLFormatter{
		data: make(map[string]interface{}),
	}
}

// YAMLString returns facts as YAML
func (formatter *YAMLFormatter) YAMLString() string {
	b, err := yaml.Marshal(formatter.data)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", b)
}

// Add puts a fact into memory for later
func (formatter *YAMLFormatter) Add(f Fact) {
	d := formatter.data
	for i, k := range f.Name {
		if i >= len(f.Name)-1 {
			d[k] = f.Value
		} else {
			newd, ok := d[k].(map[string]interface{})
			if ok {
				d = newd
			} else {
				d[k] = make(map[string]interface{})
				d = d[k].(map[string]interface{})
			}
		}
	}
}

// Finish dumps facts from memory to standard output
func (formatter *YAMLFormatter) Finish() {
	fmt.Println(formatter.YAMLString())
}
