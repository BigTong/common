package http_service

import (
	"net/http"
	"testing"
)

var handler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hehe"))
}

func TestHttps(t *testing.T) {
	service := NewService("127.0.0.1", "8080")
	service.WithHandleFunc("/hehe", handler)
	service.SupportHttps(HTTPS_ONLY, "cert.pem", "key.pem")
	service.Run()
}

func TestHttp(t *testing.T) {
	service := NewService("127.0.0.1", "8080")
	service.WithHandleFunc("/hehe", handler)
	service.Run()
}

func TestHttpAndHttps(t *testing.T) {
	service := NewService("127.0.0.1", "8080")
	service.WithHandleFunc("/hehe", handler)
	service.SupportHttps(HTTP_AND_HTTPS, "cert.pem", "key.pem")
	service.Run()
}
