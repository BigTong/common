package service

import (
	"net"
	"net/http"
)

const (
	UNKNOWN = iota
	HTTP_ONLY
	HTTPS_ONLY
	HTTP_AND_HTTPS
)

type Service struct {
	host         string
	httpPort     string
	httpsPort    string
	certFileName string
	keyFileName  string

	supportType int

	handlers map[string]http.HandlerFunc
}

func (service *Service) Run() {
	service.initHandleFunc()
	service.run()
}

func (service *Service) SetHandleFunc(pattern string,
	handlerFunc http.HandlerFunc) *Service {
	service.handlers[pattern] = handlerFunc
	return service
}

func (service *Service) SetHandleFuncs(
	handlers map[string]http.HandlerFunc) *Service {
	for k, v := range handlers {
		service.handlers[k] = v
	}
	return service
}

func (service *Service) SetSupportType(supportType int) *Service {
	service.supportType = supportType
	return service
}

func (service *Service) SetHttpPort(port string) *Service {
	service.httpPort = port
	return service
}

func (service *Service) SetHttpsConfig(port, cert, key string) *Service {
	service.httpsPort = port
	service.certFileName = cert
	service.keyFileName = key
	return service
}

func (service *Service) SetHttpsCertificate(cert, key string) *Service {
	service.certFileName = cert
	service.keyFileName = key
	return service
}

func (service *Service) initHandleFunc() {
	for k, v := range service.handlers {
		http.HandleFunc(k, v)
	}
}

func (service *Service) run() {
	switch service.supportType {
	case HTTP_ONLY:
		service.serveHttp()
	case HTTPS_ONLY:
		service.serveHttps()
	case HTTP_AND_HTTPS:
		go service.serveHttp()
		service.serveHttps()
	default:
		panic("not support protocol type")
	}
}

func (service *Service) serveHttp() {
	err := http.ListenAndServe(net.JoinHostPort(service.host, service.httpPort), nil)
	if err != nil {
		panic(err.Error())
	}
}

func (service *Service) serveHttps() {
	err := http.ListenAndServeTLS(net.JoinHostPort(service.host, service.httpsPort),
		service.certFileName, service.keyFileName, nil)
	if err != nil {
		panic(err.Error())
	}
}
