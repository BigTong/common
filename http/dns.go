package http_entity

import (
	"net"
)

type DnsClientInterface interface {
	ResolveTCPAddr(netw, addr string) (*net.TCPAddr, error)
}

var defaultDnsClientInstance *defaultDnsClient = &defaultDnsClient{}

type defaultDnsClient struct{}

// default dns client use net.ResolveTCPAddr
func (self *defaultDnsClient) ResolveTCPAddr(netw, addr string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr(netw, addr)
}
