package http_entity

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

const (
	ZERO_IPS_ERROR = "zero ips"
)

type DnsClient struct {
	remoteDnsServerIP string
	client            *dns.Client
}

func NewDnsClient(dnsServerAddr string) DnsClientInterface {
	if len(dnsServerAddr) == 0 {
		return nil
	}
	return &DnsClient{
		remoteDnsServerIP: dnsServerAddr,
		client:            new(dns.Client),
	}
}

func (self *DnsClient) ResolveTCPAddr(netw, addr string) (*net.TCPAddr, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, errors.New("illegal addr, cant get host and port")
	}

	ip := net.ParseIP(host)
	if ip != nil {
		return net.ResolveTCPAddr(netw, net.JoinHostPort(ip.String(), port))
	}

	candidateIp, err := self.getIpByHost(host)
	if err != nil {
		return nil, err
	}

	return net.ResolveTCPAddr(netw, net.JoinHostPort(candidateIp, port))
}

//TODO use custom local ip for dns look up
func (self *DnsClient) ResolveTCPAddrWithLAddr() {}

func (self *DnsClient) getIpByHost(host string) (string, error) {
	ip := ""

	msg := new(dns.Msg)
	msg.SetQuestion(host+".", dns.TypeA)
	resp, _, err := self.client.Exchange(msg, self.remoteDnsServerIP+":53")
	if err != nil {
		return ip, err
	}

	if len(resp.Answer) == 0 {
		return ip, errors.New(ZERO_IPS_ERROR)
	}

	ips := []string{}
	for _, ans := range resp.Answer {
		if ip, ok := ans.(*dns.A); ok {
			ips = append(ips, ip.A.String())
		}
	}

	length := len(ips)
	if length == 0 {
		return ip, errors.New(ZERO_IPS_ERROR)
	}

	if length == 1 {
		ip = ips[0]
	} else {
		ip = ips[rand.Intn(length)]
	}

	return ip, nil
}
