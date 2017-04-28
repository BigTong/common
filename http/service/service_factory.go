package http_service

import (
	"net/http"
)

func NewHttpService(host, port string) *Service {
	return &Service{
		host:        host,
		httpPort:    port,
		handlers:    make(map[string]http.HandlerFunc),
		supportType: HTTP_ONLY,
	}
}

func NewHttpsService(host, port string) *Service {
	return &Service{
		host:        host,
		httpsPort:   port,
		handlers:    make(map[string]http.HandlerFunc),
		supportType: HTTPS_ONLY,
	}
}

func NewHttpAndHttpsService(host, httpPort, httpsPort string) *Service {
	return &Service{
		host:        host,
		httpPort:    httpPort,
		httpsPort:   httpsPort,
		handlers:    make(map[string]http.HandlerFunc),
		supportType: HTTP_AND_HTTPS,
	}
}

func NewService(host string) *Service {
	return &Service{
		host:        host,
		handlers:    make(map[string]http.HandlerFunc),
		supportType: UNKNOWN,
	}
}
