package v1

import "net/url"

func mergeQueryValues(dst, src url.Values) url.Values {
	if len(src) == 0 {
		return dst
	}
	if dst == nil {
		dst = url.Values{}
	}
	for key, values := range src {
		for _, value := range values {
			dst.Add(key, value)
		}
	}
	return dst
}
