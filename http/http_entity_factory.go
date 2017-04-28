package http_entity

import (
	"errors"
)

// provide some static factory method for constructing http_entity

func NewHttpEntity(timeout int) (*HttpEntity, error) {
	builder := HttpEntityBuilder(timeout)
	return builder.Build(), nil
}

func NewHttpEntityWithIp(timeout int, ip string) (*HttpEntity, error) {
	if len(ip) == 0 {
		return nil, errors.New("ip address is invalid")
	}

	builder := HttpEntityBuilder(timeout)
	return builder.Ip(ip).Build(), nil
}

func NewHttpEntityWithProxy(timeout int, proxy string) (*HttpEntity, error) {
	if len(proxy) == 0 {
		return nil, errors.New("proxy address is invalid")
	}

	builder := HttpEntityBuilder(timeout)
	return builder.Proxy(proxy).Build(), nil
}
