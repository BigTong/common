package dns

import (
	"net"
)

type DnsClientInterface interface {
	ResolveTCPAddr(netw, addr string) (*net.TCPAddr, error)
}

var DefaultDnsClientInstance *DefaultDnsClient = &DefaultDnsClient{}

type DefaultDnsClient struct{}

// default dns client use net.ResolveTCPAddr
func (self *DefaultDnsClient) ResolveTCPAddr(netw, addr string) (*net.TCPAddr, error) {
	return net.ResolveTCPAddr(netw, addr)
}
