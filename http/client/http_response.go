package client

import "net/http"

// HttpResponse definition
type HttpResponse struct {
	Status  int
	Body    []byte
	Header  http.Header
	Cookies []*http.Cookie
}

func NewHttpResponse() *HttpResponse {
	return &HttpResponse{
		Status:  DefaultHttpRespStatus,
		Body:    []byte{},
		Header:  nil,
		Cookies: nil,
	}
}
