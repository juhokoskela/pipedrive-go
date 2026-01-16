package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type NumberBool bool

func (b *NumberBool) UnmarshalJSON(data []byte) error {
	if b == nil {
		return fmt.Errorf("v1.NumberBool: UnmarshalJSON on nil receiver")
	}
	data = bytes.TrimSpace(data)
	switch {
	case len(data) == 0, bytes.Equal(data, []byte("null")), bytes.Equal(data, []byte("false")), bytes.Equal(data, []byte("0")):
		*b = false
		return nil
	case bytes.Equal(data, []byte("true")), bytes.Equal(data, []byte("1")):
		*b = true
		return nil
	default:
		var numeric float64
		if err := json.Unmarshal(data, &numeric); err == nil {
			*b = numeric != 0
			return nil
		}
	}
	return fmt.Errorf("v1.NumberBool: invalid value %q", string(data))
}
