package http_service

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
	handlers     map[string]http.HandlerFunc
	supportType  int
	certFileName string
	keyFileName  string
}

func (service *Service) Run() {
	service.initHandleFunc()
	service.run()
}

func (service *Service) WithHandleFunc(pattern string,
	handlerFunc http.HandlerFunc) *Service {
	service.handlers[pattern] = handlerFunc
	return service
}

func (service *Service) WithHandleFuncs(
	handlers map[string]http.HandlerFunc) *Service {
	for k, v := range handlers {
		service.handlers[k] = v
	}
	return service
}

func (service *Service) SupportType(supportType int) *Service {
	service.supportType = supportType
	return service
}

func (service *Service) SupportHttp(port string) *Service {
	service.httpPort = port
	return service
}

func (service *Service) SupportHttps(port, cert, key string) *Service {
	service.httpsPort = port
	service.certFileName = cert
	service.keyFileName = key
	return service
}

func (service *Service) HttpsCertificate(cert, key string) *Service {
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
