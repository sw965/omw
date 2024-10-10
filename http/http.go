package http

import (
	"net/url"
	"encoding/json"
)

func QueryToStructJSON[T any](query string) (T, error) {
	decoded, err := url.QueryUnescape(query)
	var t T
	if err != nil {
		return t, err
	}
	err = json.Unmarshal([]byte(decoded), &t)
	return t, err
}