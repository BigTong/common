package client

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/BigTong/common/dns"
	"golang.org/x/net/publicsuffix"
)

// ClientBuilder pattern for construct HttpEntity
type ClientBuilder struct {
	// request timeout threshold
	timeout int

	// local ip. If this field is setted, dial remote address with this ip.
	ip string

	// support use proxy
	proxy string

	// use cookieJar
	needCookie bool

	// re-use tcp connection
	keepAlive bool

	// dns client interface, resolve host and ip address
	dnsClient dns.DnsClientInterface

	// HttpClient is product of ClientBuilder
	product *HttpClient
}

type dialFunc func(netw, addr string) (net.Conn, error)
type proxyFunc func(*http.Request) (*url.URL, error)
type redirectFunc func(req *http.Request, via []*http.Request) error

// HttpEntityClientBuilder is ClientBuilder pattern, construct a HttpClient
func HttpClientBuilder(timeout int) *ClientBuilder {
	ClientBuilder := new(ClientBuilder)
	ClientBuilder.timeout = timeout
	ClientBuilder.proxy = ""
	ClientBuilder.ip = ""
	ClientBuilder.needCookie = false
	ClientBuilder.keepAlive = false

	ClientBuilder.dnsClient = dns.DefaultDnsClientInstance
	return ClientBuilder
}

func (b *ClientBuilder) Ip(ip string) *ClientBuilder {
	if len(ip) != 0 {
		b.ip = ip
	}
	return b
}

func (b *ClientBuilder) Proxy(proxy string) *ClientBuilder {
	if len(proxy) != 0 {
		b.proxy = proxy
	}
	return b
}

func (b *ClientBuilder) DnsClient(dns dns.DnsClientInterface) *ClientBuilder {
	if dns != nil {
		b.dnsClient = dns
	}
	return b
}

func (b *ClientBuilder) CookieJar(needCookie bool) *ClientBuilder {
	b.needCookie = needCookie
	return b
}

func (b *ClientBuilder) KeepAlive(keepAlive bool) *ClientBuilder {
	b.keepAlive = keepAlive
	return b
}

func (b *ClientBuilder) Build() *HttpClient {
	b.buildHttpClient()
	return b.product
}

func (b *ClientBuilder) createHttpClient() {
	b.product = new(HttpClient)

	b.product.timeout = b.timeout
	b.product.ip = b.ip
	b.product.proxy = b.proxy
	b.product.needCookie = b.needCookie
	b.product.keepAlive = b.keepAlive

	b.product.dnsClient = b.dnsClient
}

func (b *ClientBuilder) buildHttpClient() {
	b.createHttpClient()
	client := &http.Client{
		Transport: &http.Transport{
			Dial:                  b.getDial(),
			DisableKeepAlives:     !b.keepAlive,
			ResponseHeaderTimeout: time.Duration(b.timeout) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 b.getProxy(),
		},
	}
	// client.CheckRedirect = b.getRedirect()
	client.Jar = b.getCookieJar()
	b.product.httpClient = client
}

func (b *ClientBuilder) getCookieJar() http.CookieJar {
	if b.product.needCookie {
		options := &cookiejar.Options{
			PublicSuffixList: publicsuffix.List,
		}
		jar, err := cookiejar.New(options)
		if err != nil {
			panic(err.Error())
		}
		return jar
	}
	return nil
}

/*
func (b *ClientBuilder) getRedirect() redirectFunc {
	if b.product.redirect {
		return nil
	}
	return func(req *http.Request, via []*http.Request) error {
		return errors.New("do not redirect")
	}
}
*/

// unchecked
func (b *ClientBuilder) getProxy() proxyFunc {
	if len(b.product.proxy) == 0 {
		return nil
	}
	u, err := url.Parse(b.product.proxy)
	if err != nil {
		return nil
	}
	return http.ProxyURL(u)
}

// return a dialFunc for transport dial
func (b *ClientBuilder) getDial() dialFunc {
	dial := func(netw, addr string) (net.Conn, error) {
		rAddr, err := b.product.dnsClient.ResolveTCPAddr(netw, addr)
		if err != nil {
			return nil, err
		}

		timeout := time.Duration(b.product.timeout) * time.Second
		deadline := time.Now().Add(timeout)
		if len(b.product.ip) == 0 {
			conn, err := net.DialTimeout(netw, rAddr.String(), timeout)
			if err != nil {
				return nil, err
			}
			conn.SetDeadline(deadline)
			return conn, nil
		}

		lAddr, err := net.ResolveTCPAddr(netw, b.product.ip+":0")
		if err != nil {
			return nil, err
		}

		conn, err := net.DialTCP(netw, lAddr, rAddr)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(deadline)
		return conn, nil
	}
	return dial
}
