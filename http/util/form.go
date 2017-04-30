package util

import (
	"net/url"
)

// used for urlencoded form
func UrlencodedFrom(params map[string]string) string {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	postDataStr := data.Encode()

	return postDataStr
}
