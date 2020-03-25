package ufacter

import (
	j "encoding/json"
	"fmt"
)

// JSONFormatter prints-out facts in JSON format
type JSONFormatter struct {
	data map[string]interface{}
}

// NewJSONFormatter returns new JSON formatter
func NewJSONFormatter() *JSONFormatter {
	return &JSONFormatter{
		data: make(map[string]interface{}),
	}
}

// Return facts as JSON with intent and final newline
func (jf *JSONFormatter) IntentJSONString() string {
	b, err := j.MarshalIndent(jf.data, "", "  ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s\n", b)
}

// Return facts as JSON
func (jf *JSONFormatter) JSONString() string {
	b, err := j.Marshal(jf.data)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", b)
}

func (jf *JSONFormatter) Add(f Fact) {
	d := jf.data
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

func (jf *JSONFormatter) Finish() {
	fmt.Println(jf.IntentJSONString())
}
