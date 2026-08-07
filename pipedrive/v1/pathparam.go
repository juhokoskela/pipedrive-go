package v1

import (
	"fmt"
	"math"
)

// validateID rejects identifiers that cannot be sent as a path parameter.
// Generated clients take path IDs as int, so on platforms where int is 32
// bits an int64 ID above MaxInt would silently wrap and address a
// different resource. Non-positive IDs never identify a resource.
// Once this returns nil, converting id to int is lossless.
func validateID[T ~int64](id T, label string) error {
	if id <= 0 {
		return fmt.Errorf("invalid %s %d", label, int64(id))
	}
	if int64(id) > math.MaxInt {
		return fmt.Errorf("%s %d overflows int on this platform", label, int64(id))
	}
	return nil
}

// validatePathParam rejects identifier values that URL resolution would
// collapse into a different endpoint: "" yields a trailing-slash path,
// while "." and ".." survive percent-escaping and resolve as dot segments
// when the request URL is built.
func validatePathParam(value, label string) error {
	switch value {
	case "", ".", "..":
		return fmt.Errorf("invalid %s %q", label, value)
	}
	return nil
}
