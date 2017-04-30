package service

import (
	"net/http"
	"testing"
)

var handler http.HandlerFunc = func(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("hehe"))
}

func TestHttps(t *testing.T) {
	service := NewHttpsService("127.0.0.1", "8080")
	service.SetHandleFunc("/hehe", handler)
	service.SetHttpsCertificate("cert.pem", "key.pem")
	go service.Run()
}

func TestHttp(t *testing.T) {
	service := NewHttpService("127.0.0.1", "8080")
	service.SetHandleFunc("/hehe", handler)
	go service.Run()
}

func TestHttpAndHttps(t *testing.T) {
	service := NewHttpAndHttpsService("127.0.0.1", "8080", "8081")
	service.SetHandleFunc("/hehe", handler)
	service.SetHttpsCertificate("cert.pem", "key.pem")
	go service.Run()
}
