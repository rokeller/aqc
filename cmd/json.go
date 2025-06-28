package cmd

import (
	"encoding/json"
	"fmt"
)

type jsonValueKind int

const (
	Array jsonValueKind = iota
	Object
	String
	Number
	Boolean
	Null
)

type jsonValue struct {
	kind    jsonValueKind
	array   []any
	object  map[string]any
	str     string
	num     float64
	boolean bool
}

var _ json.Unmarshaler = (*jsonValue)(nil)

func (j *jsonValue) Get() any {
	switch j.kind {
	case Null:
		return nil

	case Boolean:
		return j.boolean

	case Number:
		return j.num

	case String:
		return j.str

	case Array:
		return j.array

	case Object:
		return j.object
	}

	panic(fmt.Sprintf("unsupported jsonValueKind: %v", j.kind))
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *jsonValue) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		j.kind = Null
		return nil
	} else if string(data) == "true" {
		j.kind = Boolean
		j.boolean = true
		return nil
	} else if string(data) == "false" {
		j.kind = Boolean
		j.boolean = false
		return nil
	} else if len(data) > 0 && data[0] == '[' {
		j.kind = Array
		return json.Unmarshal(data, &j.array)
	} else if len(data) > 0 && data[0] == '{' {
		j.kind = Object
		return json.Unmarshal(data, &j.object)
	} else if len(data) > 0 && data[0] == '"' {
		j.kind = String
		return json.Unmarshal(data, &j.str)
	} else {
		j.kind = Number
		return json.Unmarshal(data, &j.num)
	}
}
