package http_entity

import (
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/publicsuffix"
)

// Builder pattern for construct HttpEntity
type Builder struct {
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
	dnsClient DnsClientInterface

	// HttpEntity is product of Builder
	product *HttpEntity
}

type dialFunc func(netw, addr string) (net.Conn, error)
type proxyFunc func(*http.Request) (*url.URL, error)
type redirectFunc func(req *http.Request, via []*http.Request) error

// HttpEntityBuilder is builder pattern, construct a HttpEntity
func HttpEntityBuilder(timeout int) *Builder {
	builder := new(Builder)
	builder.timeout = timeout
	builder.proxy = ""
	builder.ip = ""
	builder.needCookie = false
	builder.keepAlive = false

	builder.dnsClient = defaultDnsClientInstance
	return builder
}

func (b *Builder) Ip(ip string) *Builder {
	if len(ip) != 0 {
		b.ip = ip
	}
	return b
}

func (b *Builder) Proxy(proxy string) *Builder {
	if len(proxy) != 0 {
		b.proxy = proxy
	}
	return b
}

func (b *Builder) DnsClient(dns DnsClientInterface) *Builder {
	if dns != nil {
		b.dnsClient = dns
	}
	return b
}

func (b *Builder) CookieJar(needCookie bool) *Builder {
	b.needCookie = needCookie
	return b
}

func (b *Builder) KeepAlive(keepAlive bool) *Builder {
	b.keepAlive = keepAlive
	return b
}

func (b *Builder) Build() *HttpEntity {
	b.buildHttpEntity()
	return b.product
}

func (b *Builder) createHttpEntity() {
	b.product = new(HttpEntity)

	b.product.timeout = b.timeout
	b.product.ip = b.ip
	b.product.proxy = b.proxy
	b.product.needCookie = b.needCookie
	b.product.keepAlive = b.keepAlive

	b.product.dnsClient = b.dnsClient
}

func (b *Builder) buildHttpEntity() {
	b.createHttpEntity()
	client := &http.Client{
		Transport: &http.Transport{
			Dial:                  b.getDial(),
			DisableKeepAlives:     !b.keepAlive,
			ResponseHeaderTimeout: time.Duration(b.timeout) * time.Second,
			TLSClientConfig:       &tls.Config{InsecureSkipVerify: true},
			Proxy:                 b.getProxy(),
		},
	}
	//client.CheckRedirect = b.getRedirect()
	client.Jar = b.getCookieJar()
	b.product.httpClient = client
}

func (b *Builder) getCookieJar() http.CookieJar {
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
func (b *Builder) getRedirect() redirectFunc {
	if b.product.redirect {
		return nil
	}
	return func(req *http.Request, via []*http.Request) error {
		return errors.New("do not redirect")
	}
}
*/

// unchecked
func (b *Builder) getProxy() proxyFunc {
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
func (b *Builder) getDial() dialFunc {
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
