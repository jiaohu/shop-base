package models

import (
	"encoding/json"
	"fmt"
)

type StringArray []string

func (s *StringArray) Scan(value any) error {
	if value == nil {
		*s = []string{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for StringArray")
	}

	return json.Unmarshal(bytes, s)
}

type Uint64Array []uint64

func (s *Uint64Array) Scan(value any) error {
	if value == nil {
		*s = []uint64{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("invalid type for Uint64Array")
	}

	if len(bytes) == 0 {
		*s = []uint64{}
		return nil
	}

	return json.Unmarshal(bytes, s)
}

type AttrContentEntity struct {
	Name      string   `json:"name"`
	Value     []string `json:"value"`
	TempValue string   `json:"tempValue"`
}

type AttrContent []AttrContentEntity

func (a *AttrContent) Scan(value any) error {
	if value == nil {
		*a = AttrContent{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("attrcontent: unsupported type %T", value)
	}

	if len(bytes) == 0 {
		*a = AttrContent{}
		return nil
	}

	return json.Unmarshal(bytes, a)
}
