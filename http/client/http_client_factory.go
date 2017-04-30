package client

import (
	"errors"
)

// provide some static factory method for constructing http_entity

func NewHttpClient(timeout int) (*HttpClient, error) {
	builder := HttpClientBuilder(timeout)
	return builder.Build(), nil
}

func NewHttpClientWithIp(timeout int, ip string) (*HttpClient, error) {
	if len(ip) == 0 {
		return nil, errors.New("ip address is invalid")
	}

	builder := HttpClientBuilder(timeout)
	return builder.Ip(ip).Build(), nil
}

func NewHttpClientWithProxy(timeout int, proxy string) (*HttpClient, error) {
	if len(proxy) == 0 {
		return nil, errors.New("proxy address is invalid")
	}

	builder := HttpClientBuilder(timeout)
	return builder.Proxy(proxy).Build(), nil
}
