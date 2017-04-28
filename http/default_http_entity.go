package http_entity

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"common/http_entity/utils"
)

const (
	defaultTimeout = 40
)

// default http entity, without any and http config such as redirect, cookiejar
// package level method, more easily for invoking
var defaultHttpEntity *HttpEntity = HttpEntityBuilder(defaultTimeout).Build()

func GetDefaultHttpEntity() *HttpEntity {
	return defaultHttpEntity
}

// simple get method, use this method without specific option, use more easily
func SimpleGet(url string, header map[string]string) (*HttpResponse, error) {
	return defaultHttpEntity.SimpleGet(url, header)
}

// GET method, only use url encoded params
func Get(url string, headers map[string]string,
	params map[string]string, cookies []*http.Cookie,
	option *FetchOption) (*HttpResponse, error) {
	return defaultHttpEntity.Get(url, headers, params, cookies, option)
}

// simple post method, use url encoded content type as params
func UrlEncodedPost(url string, header map[string]string,
	params map[string]string) (*HttpResponse, error) {
	return defaultHttpEntity.Post(url, header,
		strings.NewReader(http_utils.UrlencodedFrom(params)), FORM_URLENCODED, nil, nil)
}

// simple post method, use json type as params
func JsonPost(url string, header map[string]string,
	body []byte) (*HttpResponse, error) {
	return defaultHttpEntity.Post(url, header, bytes.NewReader(body), JSON, nil, nil)
}

// POST method, accept body as post params, auto add Content-Type in header
// according to FormContentType
func Post(url string, headers map[string]string, body io.Reader,
	contentType FormContentType, cookies []*http.Cookie, option *FetchOption) (*HttpResponse, error) {
	return defaultHttpEntity.Post(url, headers, body, contentType, cookies, option)
}
