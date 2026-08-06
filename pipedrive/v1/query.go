package v1

import "net/url"

// mergeQueryValues merges src into dst, replacing any values dst already
// holds for a key. Options are applied in order, so the last option to set
// a key wins; appending instead would emit the key twice and leave
// precedence up to the server.
func mergeQueryValues(dst, src url.Values) url.Values {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = url.Values{}
	}
	for key, values := range src {
		dst[key] = append([]string(nil), values...)
	}
	return dst
}
