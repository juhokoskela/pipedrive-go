package v1

import (
	"bytes"
	"fmt"
	"strconv"
	"time"
)

var v1TimeLayouts = []string{
	time.RFC3339,
	v1DateTimeLayout,
}

type DateTime struct {
	time.Time
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	if t == nil {
		return fmt.Errorf("v1.DateTime: UnmarshalJSON on nil receiver")
	}
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}

	value, err := strconv.Unquote(string(data))
	if err != nil {
		return fmt.Errorf("v1.DateTime: decode: %w", err)
	}
	if value == "" {
		return nil
	}

	for _, layout := range v1TimeLayouts {
		parsed, err := time.Parse(layout, value)
		if err == nil {
			t.Time = parsed
			return nil
		}
	}

	return fmt.Errorf("v1.DateTime: parse %q", value)
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(strconv.Quote(formatV1Time(t.Time))), nil
}
